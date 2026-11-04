package services

import (
	"encoding/json"

	"agent-desk/internal/pkg/dto"
	"agent-desk/internal/pkg/dto/response"
	"agent-desk/internal/pkg/enums"
	"agent-desk/internal/pkg/errorsx"
	"agent-desk/internal/pkg/i18nx"

	"github.com/mlogclub/simple/sqls"
)

func (s *systemConfigService) GetSiteLanguage() string {
	admin, _ := s.GetLanguages()
	return admin
}

func (s *systemConfigService) GetClientLanguage() string {
	_, client := s.GetLanguages()
	return client
}

func (s *systemConfigService) GetLanguages() (admin, client string) {
	admin = i18nx.LocaleZhCN
	if sqls.DB() == nil {
		return admin, admin
	}
	var configuredClient string
	for _, item := range s.FindByGroupCode("general") {
		if item.Status == enums.StatusOk && (item.ConfigKey == "language" || item.ConfigKey == "clientLanguage") {
			var language string
			if json.Unmarshal([]byte(item.ConfigValue), &language) == nil && ((item.ConfigKey == "language" && isSiteLanguage(language)) || (item.ConfigKey == "clientLanguage" && isClientLanguage(language))) {
				if item.ConfigKey == "language" {
					admin = language
				} else {
					configuredClient = language
				}
			}
		}
	}
	if configuredClient == "" {
		configuredClient = admin
	}
	return admin, configuredClient
}

func (s *systemConfigService) SaveSiteLanguage(language string, operator *dto.AuthPrincipal) error {
	if operator == nil {
		return errorsx.UnauthorizedI18n("error.auth.expired")
	}
	raw, err := json.Marshal(language)
	if err != nil {
		return err
	}
	return s.SaveGroupConfig("general", map[string]json.RawMessage{"language": raw}, operator)
}

func (s *systemConfigService) SaveLanguages(language, clientLanguage string, operator *dto.AuthPrincipal) error {
	if operator == nil {
		return errorsx.UnauthorizedI18n("error.auth.expired")
	}
	adminRaw, _ := json.Marshal(language)
	clientRaw, _ := json.Marshal(clientLanguage)
	return s.SaveGroupConfig("general", map[string]json.RawMessage{
		"language": adminRaw, "clientLanguage": clientRaw,
	}, operator)
}

func isSiteLanguage(language string) bool {
	return language == i18nx.LocaleZhCN || language == i18nx.LocaleEnUS
}

func isClientLanguage(language string) bool {
	return isSiteLanguage(language) || language == i18nx.LocaleIdID
}

type clientLanguageValidator struct{}

func (clientLanguageValidator) Validate(raw json.RawMessage) (json.RawMessage, []response.ConfigFieldError, error) {
	var language string
	if json.Unmarshal(raw, &language) != nil || !isClientLanguage(language) {
		return nil, nil, errorsx.InvalidParamI18n("systemConfig.general.clientLanguage.invalid")
	}
	return raw, nil, nil
}

type siteLanguageValidator struct{}

func (siteLanguageValidator) Validate(raw json.RawMessage) (json.RawMessage, []response.ConfigFieldError, error) {
	var language string
	if json.Unmarshal(raw, &language) != nil || !isSiteLanguage(language) {
		return nil, nil, errorsx.InvalidParamI18n("systemConfig.general.language.invalid")
	}
	return raw, nil, nil
}
