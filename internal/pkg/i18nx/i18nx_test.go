package i18nx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNormalizeLocale(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "default for blank", in: "", want: LocaleZhCN},
		{name: "exact chinese", in: "zh-CN", want: LocaleZhCN},
		{name: "underscore chinese", in: "zh_CN", want: LocaleZhCN},
		{name: "short chinese", in: "zh", want: LocaleZhCN},
		{name: "exact english", in: "en-US", want: LocaleEnUS},
		{name: "underscore english", in: "en_US", want: LocaleEnUS},
		{name: "short english", in: "en", want: LocaleEnUS},
		{name: "indonesian", in: "id-ID", want: LocaleIdID},
		{name: "short indonesian", in: "id", want: LocaleIdID},
		{name: "unsupported falls back", in: "fr-FR", want: LocaleZhCN},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeLocale(tt.in); got != tt.want {
				t.Fatalf("NormalizeLocale(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestResolveLocaleUsesDefaultLocale(t *testing.T) {
	SetDefaultLocale(LocaleZhCN)
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/user/list", nil)
	req.Header.Set("Accept-Language", "fr-FR, en-US;q=0.9, zh-CN;q=0.8")

	if got := ResolveRequestLocale(req); got != LocaleZhCN {
		t.Fatalf("ResolveRequestLocale() = %q, want %q", got, LocaleZhCN)
	}
}

func TestResolveLocaleIgnoresRequestLocaleHeaders(t *testing.T) {
	SetDefaultLocale(LocaleZhCN)
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/user/list", nil)
	req.Header.Set("X-Locale", "en-US")
	req.Header.Set("Accept-Language", "zh-CN")

	if got := ResolveRequestLocale(req); got != LocaleZhCN {
		t.Fatalf("ResolveRequestLocale() = %q, want %q", got, LocaleZhCN)
	}
}

func TestTranslateUsesConfiguredDefaultForUnsupportedLocale(t *testing.T) {
	SetDefaultLocale(LocaleZhCN)
	if got := TLocale(LocaleIdID, "error.e0210"); got != "Saluran tidak tersedia." {
		t.Fatalf("Indonesian channel error = %q", got)
	}
	if got := TLocale(LocaleIdID, "error.e0117"); got != "The conversation does not belong to this AI Agent." {
		t.Fatalf("English fallback = %q", got)
	}
	if got := TLocale(LocaleEnUS, "error.auth.expired"); got != "Your session has expired. Please sign in again." {
		t.Fatalf("english translation = %q", got)
	}
	if got := TLocale("fr-FR", "error.auth.expired"); got != "未登录或登录已过期" {
		t.Fatalf("fallback translation = %q", got)
	}
}

func TestIndonesianClientErrorsDoNotFallBackToEnglish(t *testing.T) {
	loadMessages()
	for _, key := range []string{
		"error.e0030", "error.e0064", "error.e0079", "error.e0090",
		"error.e0116", "error.e0119", "error.e0120", "error.e0150",
		"error.e0157", "error.e0160", "error.e0161", "error.e0209",
		"error.e0211", "error.e0222", "error.e0245", "error.e0265",
		"error.e0322", "error.e0323", "error.e0324",
	} {
		if value := lookupLocale(LocaleIdID, key); value == "" {
			t.Errorf("missing Indonesian translation for %s", key)
		}
	}
}

func TestGetfSupportsLocaleKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		locale string
		key    string
		args   []any
		want   string
	}{
		{
			name:   "batch limit",
			locale: LocaleEnUS,
			key:    "error.agentTeamSchedule.batchLimit",
			args:   []any{31},
			want:   "You can generate at most 31 schedule entries at a time.",
		},
		{
			name:   "directory not found",
			locale: LocaleEnUS,
			key:    "error.e0273",
			want:   "Directory not found.",
		},
		{
			name:   "single knowledge base reference",
			locale: LocaleEnUS,
			key:    "error.knowledgeBase.referencedByAgent",
			args:   []any{"SalesBot"},
			want:   "This knowledge base is referenced by AI Agent \"SalesBot\". Remove the binding first.",
		},
		{
			name:   "multiple knowledge base references",
			locale: LocaleEnUS,
			key:    "error.knowledgeBase.referencedByAgents",
			args:   []any{3},
			want:   "This knowledge base is referenced by 3 AI Agents. Remove the bindings first.",
		},
		{
			name:   "faq import duplicated question",
			locale: LocaleEnUS,
			key:    "error.knowledgeFAQImport.duplicateQuestionInFile",
			args:   []any{12},
			want:   "The standard question is duplicated in the same file. It first appeared on row 12.",
		},
		{
			name:   "zh keeps chinese value",
			locale: LocaleZhCN,
			key:    "error.e0273",
			want:   "目录不存在",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Getf(tt.locale, tt.key, tt.args...); got != tt.want {
				t.Fatalf("Getf() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetfFallsBackToKey(t *testing.T) {
	t.Parallel()

	if got := Getf(LocaleEnUS, "missing.key"); got != "missing.key" {
		t.Fatalf("missing translation = %q, want %q", got, "missing.key")
	}
}

func TestMiddlewareStoresLocale(t *testing.T) {
	SetDefaultLocale(LocaleZhCN)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Middleware())
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, Locale(ctx))
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("X-Locale", "en-US")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Body.String(); got != LocaleZhCN {
		t.Fatalf("middleware locale = %q, want %q", got, LocaleZhCN)
	}
}
