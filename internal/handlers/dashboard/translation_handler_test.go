package dashboard

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"agent-desk/internal/pkg/constants"
	"agent-desk/internal/pkg/dto"

	"github.com/gin-gonic/gin"
)

func TestTranslationPermissions(t *testing.T) {
	for _, test := range []struct {
		name, body  string
		permissions []string
		code        int
	}{
		{"no view", `{"conversationId":1,"messageId":1,"targetLanguage":"en"}`, nil, 3001},
		{"draft needs send", `{"conversationId":1,"text":"test","format":"text","targetLanguage":"en"}`, []string{constants.PermissionConversationView.Code}, 3001},
		{"reading needs only view", `{"conversationId":1,"messageId":1,"targetLanguage":"xx"}`, []string{constants.PermissionConversationView.Code}, 1000},
		{"draft allowed to validate", `{"conversationId":1,"text":"test","targetLanguage":"xx"}`, []string{constants.PermissionConversationView.Code, constants.PermissionConversationSend.Code}, 1000},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Set("authPrincipal", &dto.AuthPrincipal{UserID: 7, Permissions: test.permissions})
			ctx.Request = httptest.NewRequest("POST", "/api/dashboard/conversation/translate", strings.NewReader(test.body))
			ctx.Request.Header.Set("Content-Type", "application/json")
			ConversationPostTranslate(ctx)
			var result struct {
				ErrorCode int `json:"errorCode"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if result.ErrorCode != test.code {
				t.Fatalf("error code = %d, want %d", result.ErrorCode, test.code)
			}
		})
	}
}
