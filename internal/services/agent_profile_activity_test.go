package services

import (
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"agent-desk/internal/models"
	"agent-desk/internal/pkg/dto"
	"agent-desk/internal/pkg/dto/request"
	"agent-desk/internal/pkg/enums"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/gorilla/websocket"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestAgentProfileCreateRecordsStatusTime(t *testing.T) {
	setupAgentActivityTestDB(t)
	before := time.Now().Add(-time.Second)
	profile, err := AgentProfileService.CreateAgentProfile(agentActivityRequest(), &dto.AuthPrincipal{UserID: 1})
	if err != nil {
		t.Fatal(err)
	}
	stored := AgentProfileService.Get(profile.ID)
	if stored.LastStatusAt == nil || stored.LastStatusAt.Before(before) {
		t.Fatalf("new profile must record its initial status time: %+v", stored)
	}
	if stored.LastOnlineAt != nil {
		t.Fatal("creating a profile must not mark the agent online")
	}
}

func TestAgentProfileUpdateRecordsOnlyStatusChanges(t *testing.T) {
	for _, status := range []enums.ServiceStatus{enums.ServiceStatusIdle, enums.ServiceStatusBusy} {
		t.Run(enums.GetServiceStatusLabel(status), func(t *testing.T) {
			db := setupAgentActivityTestDB(t)
			previous := time.Now().Add(-time.Hour).Truncate(time.Second)
			profile := models.AgentProfile{
				UserID: 1, TeamID: 1, AgentCode: "001", DisplayName: "agent",
				ServiceStatus: enums.ServiceStatusIdle, LastStatusAt: &previous,
			}
			if err := db.Create(&profile).Error; err != nil {
				t.Fatal(err)
			}
			req := agentActivityRequest()
			req.ServiceStatus = status
			req.DisplayName = "renamed"
			before := time.Now().Add(-time.Second)
			if err := AgentProfileService.UpdateAgentProfile(request.UpdateAgentProfileRequest{
				ID: profile.ID, CreateAgentProfileRequest: req,
			}, &dto.AuthPrincipal{UserID: 1}); err != nil {
				t.Fatal(err)
			}
			stored := AgentProfileService.Get(profile.ID)
			if stored.ServiceStatus != status || stored.DisplayName != req.DisplayName {
				t.Fatalf("profile update was not saved: %+v", stored)
			}
			if stored.LastStatusAt == nil {
				t.Fatal("status time must not be cleared")
			}
			if status == enums.ServiceStatusIdle && !stored.LastStatusAt.Equal(previous) {
				t.Fatal("saving an unchanged status must preserve its timestamp")
			}
			if status == enums.ServiceStatusBusy && stored.LastStatusAt.Before(before) {
				t.Fatal("changing status must update its timestamp")
			}
		})
	}
}

func TestWsAgentOnlineActivity(t *testing.T) {
	for _, tc := range []struct {
		name      string
		role      string
		principal *dto.AuthPrincipal
		online    bool
	}{
		{"dashboard", realtimeRoleAdmin, &dto.AuthPrincipal{UserID: 1, UserType: enums.UserTypeEmployee}, true},
		{"notification", realtimeRoleNotification, &dto.AuthPrincipal{UserID: 1, UserType: enums.UserTypeEmployee}, true},
		{"customer", realtimeRoleUser, &dto.AuthPrincipal{UserID: 1, UserType: enums.UserTypeEmployee}, false},
		{"non_employee", realtimeRoleAdmin, &dto.AuthPrincipal{UserID: 1, UserType: enums.UserTypeUser}, false},
		{"anonymous", realtimeRoleUser, nil, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := setupAgentActivityTestDB(t)
			previous := time.Now().Add(-time.Hour).Truncate(time.Second)
			profile := models.AgentProfile{
				UserID: 1, AgentCode: "001", LastStatusAt: &previous,
				AuditFields: models.AuditFields{UpdatedAt: previous},
			}
			if err := db.Create(&profile).Error; err != nil {
				t.Fatal(err)
			}
			svc := newWsService()
			router := gin.New()
			router.GET("/ws", func(ctx *gin.Context) {
				if err := svc.upgradeConnection(ctx, tc.principal, nil, tc.role); err != nil {
					t.Errorf("upgrade websocket: %v", err)
				}
			})
			server := httptest.NewServer(router)
			defer server.Close()
			dial := func() *websocket.Conn {
				t.Helper()
				conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/ws", nil)
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = conn.Close() })
				readAgentActivityEvent(t, conn, enums.IMRealtimeEventConnected)
				return conn
			}
			first := dial()
			stored := AgentProfileService.Get(profile.ID)
			if !tc.online {
				if stored.LastOnlineAt != nil {
					t.Fatal("customer connections must not mark an agent online")
				}
				return
			}
			if stored.LastOnlineAt == nil || stored.LastOnlineAt.Before(time.Now().Add(-time.Minute)) {
				t.Fatal("staff websocket connection must record online time")
			}
			second := dial()
			if err := first.Close(); err != nil {
				t.Fatal(err)
			}
			for _, lastOnline := range []time.Time{previous, time.Now().Add(-time.Second).Truncate(time.Second)} {
				if err := db.Model(&profile).UpdateColumn("last_online_at", lastOnline).Error; err != nil {
					t.Fatal(err)
				}
				if err := second.WriteJSON(realtimeClientMessage{Type: enums.IMRealtimeClientTypePing}); err != nil {
					t.Fatal(err)
				}
				readAgentActivityEvent(t, second, enums.IMRealtimeEventPong)
				stored = AgentProfileService.Get(profile.ID)
				if stored.LastOnlineAt == nil {
					t.Fatal("heartbeat must retain an online time")
				}
				if lastOnline.Equal(previous) && !stored.LastOnlineAt.After(previous) {
					t.Fatal("remaining tab's heartbeat must refresh stale online time")
				}
				if !lastOnline.Equal(previous) && !stored.LastOnlineAt.Equal(lastOnline) {
					t.Fatal("heartbeats within 30 seconds must not rewrite online time")
				}
			}
			if !stored.UpdatedAt.Equal(previous) || !stored.LastStatusAt.Equal(previous) {
				t.Fatal("online heartbeat must not change profile audit or service status time")
			}
		})
	}
}

func readAgentActivityEvent(t *testing.T, conn *websocket.Conn, eventType string) {
	t.Helper()
	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	for {
		var event struct {
			Type string `json:"type"`
		}
		if err := conn.ReadJSON(&event); err != nil {
			t.Fatal(err)
		}
		if event.Type == eventType {
			return
		}
	}
}

func setupAgentActivityTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "activity.db")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "t_", SingularTable: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db.AutoMigrate(&models.User{}, &models.AgentTeam{}, &models.AgentProfile{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.User{ID: 1, Username: "agent", UserType: enums.UserTypeEmployee}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.AgentTeam{ID: 1, Name: "support"}).Error; err != nil {
		t.Fatal(err)
	}
	sqls.SetDB(db)
	return db
}

func agentActivityRequest() request.CreateAgentProfileRequest {
	return request.CreateAgentProfileRequest{
		UserID: 1, TeamID: 1, AgentCode: "001", DisplayName: "agent", ServiceStatus: enums.ServiceStatusIdle,
	}
}
