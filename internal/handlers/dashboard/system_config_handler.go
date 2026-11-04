package dashboard

import (
	"agent-desk/internal/pkg/constants"
	"agent-desk/internal/pkg/dto/request"
	"agent-desk/internal/pkg/errorsx"
	"agent-desk/internal/pkg/httpx"
	"agent-desk/internal/pkg/httpx/params"
	"agent-desk/internal/services"

	"github.com/gin-gonic/gin"
)

func SystemConfigPostLanguage(ctx *gin.Context) {
	operator, err := services.AuthService.RequirePermission(ctx, constants.PermissionSystemConfigUpdate)
	if err != nil {
		httpx.WriteJSON(ctx, err)
		return
	}
	var req request.UpdateSystemLanguageRequest
	if err := params.ReadJSON(ctx, &req); err != nil {
		httpx.WriteJSON(ctx, errorsx.InvalidParamI18n("systemConfig.general.language.invalid"))
		return
	}
	httpx.WriteJSON(ctx, services.SystemConfigService.SaveLanguages(req.Language, req.ClientLanguage, operator))
}
