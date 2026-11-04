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

func TestSystemLanguagePermissionAndValidation(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		permission bool
		code       int
	}{
		{"forbidden", `{"language":"en-US","clientLanguage":"zh-CN"}`, false, 3001},
		{"invalid language", `{"language":"id","clientLanguage":"zh-CN"}`, true, 1000},
		{"invalid client language", `{"language":"en-US","clientLanguage":"id"}`, true, 1000},
		{"missing client language", `{"language":"en-US"}`, true, 1000},
		{"invalid format", `{"language":123,"clientLanguage":"zh-CN"}`, true, 1000},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			principal := &dto.AuthPrincipal{UserID: 1}
			if tc.permission {
				principal.Permissions = []string{constants.PermissionSystemConfigUpdate.Code}
			}
			ctx.Set("authPrincipal", principal)
			ctx.Request = httptest.NewRequest("POST", "/api/dashboard/system-config/language", strings.NewReader(tc.body))
			ctx.Request.Header.Set("Content-Type", "application/json")
			SystemConfigPostLanguage(ctx)
			var result struct {
				ErrorCode int `json:"errorCode"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if result.ErrorCode != tc.code {
				t.Fatalf("error code = %d, want %d", result.ErrorCode, tc.code)
			}
		})
	}
}
