package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"agent-desk/internal/ai"
	"agent-desk/internal/models"
	"agent-desk/internal/pkg/dto"
	"agent-desk/internal/pkg/dto/request"
	"agent-desk/internal/pkg/enums"
	"agent-desk/internal/pkg/errorsx"
)

func TestTranslationDocument(t *testing.T) {
	source := `<p>请查看 UID 100003，订单号 ABC-900，金额 Rp 10.000。</p><img data-asset-id="asset-1" data-provider="local" data-storage-key="test/a.png" alt="receipt"><p><a href="https://example.com/a?x=1">打开链接</a> https://example.com/order/99</p>`
	doc, err := prepareTranslation(source, "html")
	if err != nil {
		t.Fatal(err)
	}
	segments := append([]translationSegment(nil), doc.segments...)
	segments[0].Text = strings.ReplaceAll(segments[0].Text, "请查看", "Please check")
	segments[1].Text = "Open link"
	payload, _ := json.Marshal(map[string]any{"segments": segments})
	result, err := doc.apply(string(payload))
	if err != nil {
		t.Fatal(err)
	}
	for _, preserved := range []string{"UID 100003", "ABC-900", "Rp 10.000", `data-asset-id="asset-1"`, `data-storage-key="test/a.png"`, `alt="receipt"`, `href="https://example.com/a?x=1"`, "https://example.com/order/99"} {
		if !strings.Contains(result.TranslatedHTML, preserved) {
			t.Fatalf("missing preserved value %s: %s", preserved, result.TranslatedHTML)
		}
	}
	if strings.Index(result.TranslatedHTML, "Please check") > strings.Index(result.TranslatedHTML, "<img") || strings.Index(result.TranslatedHTML, "<img") > strings.Index(result.TranslatedHTML, "Open link") {
		t.Fatal("image order changed")
	}
	if !strings.Contains(result.TranslatedText, "\n") {
		t.Fatal("paragraph structure lost")
	}
	for _, output := range []string{
		`not json`, `{"segments":[]}`, `{"segments":[{"text":"test"}]}`,
		`{"segments":[{"id":0,"text":"test","html":"<img>"}]}`,
		`{"segments":[{"id":0,"text":"no identifier"}]}`,
		`{"segments":[{"id":0,"text":"__TR_KEEP_0__ __TR_KEEP_0__"}]}`,
		`{"segments":[{"id":0,"text":"__TR_KEEP_0__ __TR_KEEP_99__"}]}`,
		`{"segments":[{"id":0,"text":"__TR_KEEP_0__"}]} {}`,
	} {
		fresh, err := prepareTranslation("UID 100003", "text")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = fresh.apply(output); err == nil {
			t.Errorf("accepted invalid output: %s", output)
		}
	}
	for _, input := range []struct{ text, format string }{
		{"", "text"}, {"  ", "text"}, {strings.Repeat("中", 4001), "text"}, {`<img src="https://example.com/a.png">`, "html"}, {"hello", "xml"}, {"__TR_KEEP_0__", "text"},
	} {
		if _, err := prepareTranslation(input.text, input.format); err == nil {
			t.Errorf("accepted invalid content: %q", input.format)
		}
	}
	if _, err := prepareTranslation(strings.Repeat("中", 4000), "text"); err != nil {
		t.Fatalf("4000 visible characters rejected: %v", err)
	}
	plain, err := prepareTranslation("Hello UID 100003\n\nOrder PW-900, amount Rp 10.000.", "text")
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(map[string]any{"segments": plain.segments})
	unchanged, err := plain.apply(string(encoded))
	if err != nil || unchanged.TranslatedText != "Hello UID 100003\n\nOrder PW-900, amount Rp 10.000." {
		t.Fatalf("line breaks / punctuation not preserved: %+v, %v", unchanged, err)
	}
}

func TestTranslationService(t *testing.T) {
	db := setupMessageWelcomeTestDB(t)
	if err := db.AutoMigrate(&models.AIConfig{}, &models.AgentRevision{}); err != nil {
		t.Fatal(err)
	}
	config := &models.AIConfig{Name: "translation-test", Status: enums.StatusOk, Provider: enums.AIProviderOpenAI, ModelType: enums.AIModelTypeLLM, ModelName: "draft-model"}
	if err := db.Create(config).Error; err != nil {
		t.Fatal(err)
	}
	translationConfig := &models.AIConfig{
		Name: "independent-translation", Status: enums.StatusDisabled,
		Provider: enums.AIProviderOpenAI, ModelType: enums.AIModelType("translation"), ModelName: "translation-model",
	}
	if err := db.Create(translationConfig).Error; err != nil {
		t.Fatal(err)
	}
	agent := createWelcomeTestAIAgent(t, db, "")
	definition, _ := json.Marshal(agentRevisionDefinition{Agent: agentRevisionAgent{AIConfigID: config.ID}, Model: agentRevisionModel{ConfigID: config.ID, Provider: string(config.Provider), ModelType: string(config.ModelType), ModelName: "published-model"}})
	revision := &models.AgentRevision{AgentID: agent.ID, Revision: 1, Status: enums.StatusOk, Definition: string(definition)}
	if err := db.Create(revision).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(agent).Updates(map[string]any{"ai_config_id": config.ID, "published_revision_id": revision.ID}).Error; err != nil {
		t.Fatal(err)
	}
	conversation := createMessageTestConversation(t, db, agent.ID)
	if err := db.Model(conversation).Updates(map[string]any{"status": enums.IMConversationStatusActive, "current_assignee_id": 7, "agent_unread_count": 3}).Error; err != nil {
		t.Fatal(err)
	}
	message := &models.Message{ConversationID: conversation.ID, SenderType: enums.IMSenderTypeCustomer, MessageType: enums.IMMessageTypeText, Content: "转人工 UID 100003", SendStatus: enums.IMMessageStatusSent}
	if err := db.Create(message).Error; err != nil {
		t.Fatal(err)
	}
	calls := 0
	var usedConfig models.AIConfig
	service := &translationService{complete: func(ctx context.Context, got models.AIConfig, system, input string) (*ai.ChatCompletionResult, error) {
		calls++
		usedConfig = got
		deadline, ok := ctx.Deadline()
		if !ok || time.Until(deadline) > 20*time.Second || got.TimeoutMS != 20000 || got.MaxRetryCount != 0 {
			t.Fatal("translation timeout not applied")
		}
		if system != translationSystemPrompt || strings.Contains(input, "100003") {
			t.Fatal("translation prompt or identifier protection not applied")
		}
		return &ai.ChatCompletionResult{Content: `{"segments":[{"id":0,"text":"Transfer to a human __TR_KEEP_0__"}]}`}, nil
	}}
	req := request.TranslateConversationRequest{ConversationID: conversation.ID, MessageID: message.ID, TargetLanguage: "en"}
	operator := &dto.AuthPrincipal{UserID: 7}
	_, err := service.Translate(context.Background(), req, operator)
	var localized *errorsx.I18nError
	if !errors.As(err, &localized) || localized.Key != "translation.modelUnavailable" || calls != 0 {
		t.Fatalf("disabled translation model must not fall back to published customer model: calls=%d, err=%v", calls, err)
	}
	if err := AIConfigService.UpdateStatus(translationConfig.ID, enums.StatusOk, operator); err != nil {
		t.Fatal(err)
	}
	result, err := service.Translate(context.Background(), req, operator)
	if err != nil || result.TranslatedText != "Transfer to a human UID 100003" {
		t.Fatalf("translation result: %+v, %v", result, err)
	}
	if usedConfig.ID != translationConfig.ID || usedConfig.ModelName != "translation-model" ||
		usedConfig.ModelType != enums.AIModelType("translation") {
		t.Fatal("translation did not use the independently configured model")
	}
	if err := db.Model(agent).Updates(map[string]any{"ai_config_id": 0, "published_revision_id": 0}).Error; err != nil {
		t.Fatal(err)
	}
	var unchanged models.Conversation
	if err := db.First(&unchanged, conversation.ID).Error; err != nil {
		t.Fatal(err)
	}
	if unchanged.Status != enums.IMConversationStatusActive || unchanged.AgentUnreadCount != 3 || unchanged.HandoffAt != nil {
		t.Fatal("translation mutated conversation / read state")
	}
	draft := "转人工 UID 100003"
	req.MessageID, req.Text = 0, &draft
	if _, err = service.Translate(context.Background(), req, operator); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []request.TranslateConversationRequest{
		{ConversationID: conversation.ID, MessageID: message.ID, TargetLanguage: "xx"},
		{ConversationID: conversation.ID, MessageID: message.ID, Text: &draft, TargetLanguage: "en"},
		{ConversationID: conversation.ID, TargetLanguage: "en"},
		{ConversationID: conversation.ID, MessageID: -1, TargetLanguage: "en"},
	} {
		if _, err = service.Translate(context.Background(), invalid, operator); err == nil {
			t.Fatal("accepted invalid request")
		}
	}
	if _, err = service.Translate(context.Background(), req, &dto.AuthPrincipal{UserID: 8}); err == nil {
		t.Fatal("other agent could translate draft")
	}
	if calls != 2 {
		t.Fatalf("invalid request called model: %d", calls)
	}
	other := createMessageTestConversation(t, db, agent.ID)
	req = request.TranslateConversationRequest{ConversationID: other.ID, MessageID: message.ID, TargetLanguage: "en"}
	if _, err = service.Translate(context.Background(), req, operator); err == nil {
		t.Fatal("cross-conversation message accepted")
	}
	req.ConversationID = conversation.ID
	if err := db.Model(message).Update("send_status", enums.IMMessageStatusRecalled).Error; err != nil {
		t.Fatal(err)
	}
	if _, err = service.Translate(context.Background(), req, operator); err == nil {
		t.Fatal("recalled message accepted")
	}
	if err := db.Model(message).Update("send_status", enums.IMMessageStatusSent).Error; err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	service.complete = func(ctx context.Context, _ models.AIConfig, _, _ string) (*ai.ChatCompletionResult, error) {
		return nil, ctx.Err()
	}
	if _, err = service.Translate(ctx, req, operator); err == nil {
		t.Fatal("cancelled request accepted")
	}
	if !errors.As(err, &localized) || localized.Key != "translation.timeout" {
		t.Fatalf("wrong cancellation error: %v", err)
	}
	deadlineCtx, stop := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer stop()
	_, err = service.Translate(deadlineCtx, req, operator)
	if !errors.As(err, &localized) || localized.Key != "translation.timeout" {
		t.Fatalf("wrong deadline error: %v", err)
	}
	service.complete = func(context.Context, models.AIConfig, string, string) (*ai.ChatCompletionResult, error) {
		return &ai.ChatCompletionResult{Content: "answer instead of JSON"}, nil
	}
	if _, err = service.Translate(context.Background(), req, operator); err == nil {
		t.Fatal("invalid model response accepted")
	}
	failedCalls := 0
	service.complete = func(context.Context, models.AIConfig, string, string) (*ai.ChatCompletionResult, error) {
		failedCalls++
		return nil, errors.New("synthetic translation provider failure")
	}
	_, err = service.Translate(context.Background(), req, operator)
	if !errors.As(err, &localized) || localized.Key != "translation.failed" || failedCalls != 1 {
		t.Fatalf("translation failure must not retry another model: calls=%d, err=%v", failedCalls, err)
	}
	if err := db.Model(conversation).Update("status", enums.IMConversationStatusClosed).Error; err != nil {
		t.Fatal(err)
	}
	req.MessageID, req.Text = 0, &draft
	if _, err = service.Translate(context.Background(), req, operator); err == nil {
		t.Fatal("closed conversation draft accepted")
	}
	if err := AIConfigService.UpdateStatus(translationConfig.ID, enums.StatusDisabled, operator); err != nil {
		t.Fatal(err)
	}
	req.MessageID, req.Text = message.ID, nil
	_, err = service.Translate(context.Background(), req, operator)
	if !errors.As(err, &localized) || localized.Key != "translation.modelUnavailable" || failedCalls != 1 {
		t.Fatal("disabled translation model called another model")
	}
	if err := db.Delete(translationConfig).Error; err != nil {
		t.Fatal(err)
	}
	_, err = service.Translate(context.Background(), req, operator)
	if !errors.As(err, &localized) || localized.Key != "translation.modelUnavailable" || failedCalls != 1 {
		t.Fatal("missing translation model called another model")
	}
}
