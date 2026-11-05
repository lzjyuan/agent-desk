package api

import (
	"testing"

	"agent-desk/internal/pkg/enums"
)

func TestResolveWidgetConfigIncludesPublicBaseURL(t *testing.T) {
	cfg, err := resolveWidgetConfig(enums.ChannelTypeWeb, `{"publicBaseUrl":"https://chat.example.com/"}`)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PublicBaseURL != "https://chat.example.com/" {
		t.Fatalf("publicBaseUrl = %q", cfg.PublicBaseURL)
	}
}
