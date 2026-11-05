package dashboard

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"agent-desk/internal/models"
	"agent-desk/internal/pkg/constants"
	"agent-desk/internal/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

func TestBuildConversationResponsesIncludesChannelName(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:conversation_channel_names?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get db: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db.AutoMigrate(&models.Channel{}, &models.ConversationReadState{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	sqls.SetDB(db)
	if err := db.Create(&models.Channel{ID: 2, Name: "slotSweet3"}).Error; err != nil {
		t.Fatalf("create channel: %v", err)
	}

	results := buildConversationResponses([]models.Conversation{
		{ID: 1, ChannelID: 2},
		{ID: 2, ChannelID: 99},
	}, "zh-CN")
	if len(results) != 2 || results[0].ChannelName != "slotSweet3" || results[1].ChannelName != "" {
		t.Fatalf("unexpected channel names: %+v", results)
	}
}

func TestConversationChannelFilterAndOptions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:conversation_channel_filter?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get db: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db.AutoMigrate(&models.Channel{}, &models.Conversation{}, &models.ConversationReadState{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	sqls.SetDB(db)
	if err := db.Create(&[]models.Channel{{ID: 2, Name: "Support", ChannelID: "support-2", ConfigJSON: `{"secret":"not-for-monitoring"}`}, {ID: 3, Name: "Support", ChannelID: "support-3"}}).Error; err != nil {
		t.Fatalf("create channels: %v", err)
	}
	if err := db.Create(&[]models.Conversation{{ID: 1, ChannelID: 2}, {ID: 2, ChannelID: 2}, {ID: 3, ChannelID: 3}}).Error; err != nil {
		t.Fatalf("create conversations: %v", err)
	}

	call := func(path string, handler gin.HandlerFunc) []byte {
		t.Helper()
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, path, nil)
		ctx.Set("authPrincipal", &dto.AuthPrincipal{Permissions: []string{constants.PermissionConversationView.Code}})
		handler(ctx)
		if recorder.Code != http.StatusOK {
			t.Fatalf("%s: status %d: %s", path, recorder.Code, recorder.Body.String())
		}
		return recorder.Body.Bytes()
	}

	var options struct {
		Data []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	optionBody := call("/api/dashboard/conversation/channel_options", ConversationGetChannel_options)
	if bytes.Contains(optionBody, []byte("not-for-monitoring")) {
		t.Fatal("channel options exposed channel config")
	}
	if err := json.Unmarshal(optionBody, &options); err != nil {
		t.Fatalf("decode options: %v", err)
	}
	if len(options.Data) != 2 || options.Data[0].ID == options.Data[1].ID {
		t.Fatalf("channel options must distinguish duplicate names: %+v", options.Data)
	}

	var page struct {
		Data struct {
			Results []struct {
				ChannelID int64 `json:"channelId"`
			} `json:"results"`
			Page struct {
				Total int64 `json:"total"`
			} `json:"page"`
		} `json:"data"`
	}
	if err := json.Unmarshal(call("/api/dashboard/conversation/list?channelId=2&page=1&limit=1", ConversationAnyList), &page); err != nil {
		t.Fatalf("decode page: %v", err)
	}
	if page.Data.Page.Total != 2 || len(page.Data.Results) != 1 || page.Data.Results[0].ChannelID != 2 {
		t.Fatalf("filter must apply before paging: %+v", page.Data)
	}
}
