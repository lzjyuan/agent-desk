package bootstrap

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"agent-desk/internal/models"
	"agent-desk/internal/pkg/config"
	"agent-desk/internal/pkg/dto"
	"agent-desk/internal/services"

	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

func TestSystemLanguageChangesWithoutRestart(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:language_%p?mode=memory&cache=shared", t)), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	previousDB := sqls.DB()
	previousConfig := config.GetCurrent()
	sqls.SetDB(db)
	t.Cleanup(func() {
		sqls.SetDB(previousDB)
		config.SetCurrent(previousConfig)
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})
	if err := db.AutoMigrate(&models.SystemConfig{}); err != nil {
		t.Fatal(err)
	}
	config.SetCurrent(&config.Config{Storage: config.StorageConfig{
		Local: config.LocalStorageConfig{Root: t.TempDir(), BaseURL: "/storage"},
	}})
	app, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{{"en-US", "zh-CN"}, {"zh-CN", "en-US"}, {"en-US", "id-ID"}} {
		language, clientLanguage := pair[0], pair[1]
		if err := services.SystemConfigService.SaveLanguages(language, clientLanguage, &dto.AuthPrincipal{UserID: 1}); err != nil {
			t.Fatal(err)
		}
		for _, path := range []string{"/api/config", "/api/dashboard/user/list", "/api/message/list"} {
			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
			var result struct {
				Message string `json:"message"`
				Data    struct {
					Language       string `json:"language"`
					ClientLanguage string `json:"clientLanguage"`
				} `json:"data"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if path == "/api/config" && result.Data.Language != language {
				t.Fatalf("public language = %s, want %s", result.Data.Language, language)
			}
			if path == "/api/config" && result.Data.ClientLanguage != clientLanguage {
				t.Fatalf("public client language = %s, want %s", result.Data.ClientLanguage, clientLanguage)
			}
			if path != "/api/config" {
				want := "未登录或登录已过期"
				responseLanguage := language
				if path == "/api/message/list" {
					responseLanguage = clientLanguage
					want = "接入渠道异常"
				}
				if responseLanguage == "en-US" {
					want = "Your session has expired. Please sign in again."
					if path == "/api/message/list" {
						want = "Channel error."
					}
				}
				if responseLanguage == "id-ID" {
					want = "Saluran tidak tersedia."
				}
				if result.Message != want {
					t.Fatalf("request language %s: message = %q, want %q", language, result.Message, want)
				}
			}
		}
	}
}
