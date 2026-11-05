package services

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestCustomerSessionExpiryIncludesTimezone(t *testing.T) {
	expiresAt := time.Date(2026, 11, 5, 7, 19, 0, 0, time.FixedZone("Saudi", 3*60*60))
	formatted := formatCustomerSessionExpiry(expiresAt)
	parsed, err := time.Parse(time.RFC3339, formatted)
	if err != nil || !parsed.Equal(expiresAt) {
		t.Fatalf("expiry %q does not preserve instant %s: %v", formatted, expiresAt, err)
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	newCustomerSessionService().SetRefreshHeaders(ctx, &CustomerSessionVerifyResult{
		Token:     "token",
		ExpiresAt: expiresAt,
		Refreshed: true,
	})
	if got := recorder.Header().Get(customerSessionExpHeader); got != formatted {
		t.Fatalf("refresh header = %q, want %q", got, formatted)
	}
}
