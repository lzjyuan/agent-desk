package request

type TranslateConversationRequest struct {
	ConversationID int64   `json:"conversationId"`
	MessageID      int64   `json:"messageId"`
	Text           *string `json:"text"`
	Format         string  `json:"format"`
	TargetLanguage string  `json:"targetLanguage"`
}
