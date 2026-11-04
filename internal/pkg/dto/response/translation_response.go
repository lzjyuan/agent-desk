package response

type TranslationResponse struct {
	TranslatedText string `json:"translatedText"`
	TranslatedHTML string `json:"translatedHTML,omitempty"`
	TargetLanguage string `json:"targetLanguage"`
}
