package services

import (
	"encoding/json"
	"testing"

	"agent-desk/internal/models"
	"agent-desk/internal/pkg/dto"
)

func TestSaveSystemLanguage(t *testing.T) {
	db := setupSystemConfigServiceTestDB(t)
	operator := &dto.AuthPrincipal{UserID: 1, Username: "admin"}
	for _, language := range []string{"en-US", "zh-CN"} {
		raw, _ := json.Marshal(language)
		if err := SystemConfigService.SaveGroupConfig("general", map[string]json.RawMessage{"language": raw}, operator); err != nil {
			t.Fatalf("save %s: %v", language, err)
		}
		item := SystemConfigService.GetByGroupAndKey("general", "language")
		if item == nil || item.ConfigValue != string(raw) || item.UpdateUserID != operator.UserID {
			t.Fatalf("language must be persisted with audit fields: %+v", item)
		}
		if got := newSystemConfigService().GetSiteLanguage(); got != language {
			t.Fatalf("new service instance language = %s, want persisted %s", got, language)
		}
	}
	for _, raw := range []string{`"fr-FR"`, `""`, `null`, `true`, `{}`} {
		if err := SystemConfigService.SaveGroupConfig("general", map[string]json.RawMessage{"language": json.RawMessage(raw)}, operator); err == nil {
			t.Fatalf("accepted invalid language %s", raw)
		}
	}
	var count int64
	if err := db.Model(&models.SystemConfig{}).Count(&count).Error; err != nil || count != 1 {
		t.Fatalf("language changes must reuse one config row: count=%d err=%v", count, err)
	}
}

func TestSystemLanguageDefaultsAndAuthentication(t *testing.T) {
	db := setupSystemConfigServiceTestDB(t)
	if got := SystemConfigService.GetSiteLanguage(); got != "zh-CN" {
		t.Fatalf("initial language = %s, want zh-CN", got)
	}
	if got := SystemConfigService.GetClientLanguage(); got != "en-US" {
		t.Fatalf("initial client language = %s, want en-US", got)
	}
	if err := SystemConfigService.SaveSiteLanguage("en-US", nil); err == nil {
		t.Fatal("language save without an operator must fail")
	}
	if err := db.Create(&models.SystemConfig{GroupCode: "general", ConfigKey: "language", ConfigValue: `"fr-FR"`}).Error; err != nil {
		t.Fatal(err)
	}
	if got := SystemConfigService.GetSiteLanguage(); got != "zh-CN" {
		t.Fatalf("corrupt language must fall back to Chinese, got %s", got)
	}
}

func TestAdminAndClientLanguagesAreIndependent(t *testing.T) {
	setupSystemConfigServiceTestDB(t)
	operator := &dto.AuthPrincipal{UserID: 1, Username: "admin"}
	if err := SystemConfigService.SaveSiteLanguage("en-US", operator); err != nil {
		t.Fatal(err)
	}
	if got := SystemConfigService.GetClientLanguage(); got != "en-US" {
		t.Fatalf("unset client language = %s, want English default", got)
	}
	if err := SystemConfigService.SaveLanguages("zh-CN", "en-US", operator); err != nil {
		t.Fatal(err)
	}
	if got := SystemConfigService.GetSiteLanguage(); got != "zh-CN" {
		t.Fatalf("admin language = %s", got)
	}
	if got := SystemConfigService.GetClientLanguage(); got != "en-US" {
		t.Fatalf("client language = %s", got)
	}
	if err := SystemConfigService.SaveLanguages("en-US", "id-ID", operator); err != nil {
		t.Fatal(err)
	}
	if got := SystemConfigService.GetClientLanguage(); got != "id-ID" {
		t.Fatalf("Indonesian client language = %s", got)
	}
	if err := SystemConfigService.SaveLanguages("id-ID", "zh-CN", operator); err == nil {
		t.Fatal("accepted Indonesian admin language")
	}
	if err := SystemConfigService.SaveLanguages("en-US", "fr-FR", operator); err == nil {
		t.Fatal("accepted unsupported client language")
	}
	if got := SystemConfigService.GetSiteLanguage(); got != "en-US" {
		t.Fatalf("failed save changed admin language to %s", got)
	}
}
