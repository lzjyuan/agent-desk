package services

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"agent-desk/internal/models"
	"agent-desk/internal/pkg/dto"
	"agent-desk/internal/pkg/dto/request"
	"agent-desk/internal/pkg/enums"

	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Opt-in smoke check: only synthetic text is sent and the database is opened read-only.
func TestTranslationLiveModel(t *testing.T) {
	path := os.Getenv("TRANSLATION_SMOKE_DB")
	if path == "" {
		t.Skip("set TRANSLATION_SMOKE_DB and TRANSLATION_SMOKE_CONVERSATION to test the configured model")
	}
	id, err := strconv.ParseInt(os.Getenv("TRANSLATION_SMOKE_CONVERSATION"), 10, 64)
	if err != nil {
		t.Fatal("TRANSLATION_SMOKE_CONVERSATION must be a valid active conversation id")
	}
	db, err := gorm.Open(sqlite.Open("file:"+filepath.ToSlash(path)+"?mode=ro"), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "t_", SingularTable: true}})
	if err != nil {
		t.Fatal("cannot open smoke database read-only")
	}
	connection, err := db.DB()
	if err != nil {
		t.Fatal("cannot access smoke database")
	}
	t.Cleanup(func() { _ = connection.Close() })
	sqls.SetDB(db)
	var conversation models.Conversation
	if err := db.First(&conversation, id).Error; err != nil || conversation.Status != enums.IMConversationStatusActive {
		t.Fatal("smoke conversation must be active")
	}
	for _, test := range []struct{ name, text, target, format string }{
		{"Chinese to English", "请核对 UID 100003，订单号 PW-900，金额 Rp 10.000。", "en", "text"},
		{"English to Indonesian", "Please check UID 100003, order PW-900, amount Rp 10.000.", "id", "text"},
		{"Indonesian to Chinese", "Mohon periksa UID 100003, nomor pesanan PW-900, jumlah Rp 10.000.", "zh-CN", "text"},
		{"handoff text only", "转人工", "id", "text"},
		{"mixed HTML", `<p>Mohon tunggu UID 100003.</p><img data-asset-id="smoke-image" data-provider="local" data-storage-key="smoke.png"><p><a href="https://example.com/order/PW-900">Lihat pesanan PW-900</a></p>`, "en", "html"},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := TranslationService.Translate(context.Background(), request.TranslateConversationRequest{ConversationID: id, Text: &test.text, TargetLanguage: test.target, Format: test.format}, &dto.AuthPrincipal{UserID: conversation.CurrentAssigneeID})
			if err != nil {
				t.Fatal(err)
			}
			if result.TranslatedText == "" || strings.Contains(result.TranslatedText, "__TR_KEEP_") {
				t.Fatal("empty or unrestored translation")
			}
			if strings.Contains(test.text, "100003") && !strings.Contains(result.TranslatedText, "100003") {
				t.Fatal("UID changed")
			}
			if test.format == "html" && (!strings.Contains(result.TranslatedHTML, `data-asset-id="smoke-image"`) || !strings.Contains(result.TranslatedHTML, `href="https://example.com/order/PW-900"`)) {
				t.Fatal("HTML assets or URL changed")
			}
			t.Logf("synthetic translation: %s", result.TranslatedText)
		})
	}
}
