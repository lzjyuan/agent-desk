package request

type UpdateSystemLanguageRequest struct {
	Language       string `json:"language"`
	ClientLanguage string `json:"clientLanguage"`
}
