package services

import (
	"encoding/json"
	"testing"
)

func TestParseWebChannelPublicBaseURL(t *testing.T) {
	for _, tc := range []struct {
		url   string
		valid bool
	}{
		{"", false},
		{"https://chat.honertop.com/", true},
		{"http://localhost:8083", true},
		{"https://chat.honertop.com/support/chat/", false},
		{"javascript:alert(1)", false},
		{"https://user:pass@chat.honertop.com/", false},
	} {
		raw, err := json.Marshal(map[string]string{"publicBaseUrl": tc.url})
		if err != nil {
			t.Fatal(err)
		}
		cfg, err := newChannelService().ParseWebChannelConfig(string(raw))
		if (err == nil) != tc.valid {
			t.Errorf("url %q: error = %v, want valid = %t", tc.url, err, tc.valid)
		}
		if tc.valid && cfg.PublicBaseURL != tc.url {
			t.Errorf("url %q: parsed = %q", tc.url, cfg.PublicBaseURL)
		}
	}
	if _, err := newChannelService().ParseWebChannelConfig(""); err == nil {
		t.Error("missing publicBaseUrl must be rejected")
	}
}
