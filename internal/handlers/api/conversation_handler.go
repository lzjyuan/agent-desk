package api

import (
	"agent-desk/internal/builders"
	"agent-desk/internal/pkg/dto/request"
	"agent-desk/internal/pkg/dto/response"
	"agent-desk/internal/pkg/httpx"
	"agent-desk/internal/pkg/i18nx"
	"agent-desk/internal/services"

	"agent-desk/internal/pkg/httpx/params"

	"github.com/gin-gonic/gin"
)

func ConversationGetBy(ctx *gin.Context) {
	id, ok := httpx.GetPathInt64(ctx, "id")
	if !ok {
		return
	}
	if services.ChannelService.GetEnabledChannel(ctx) == nil {
		httpx.WriteJSON(ctx, httpx.JsonErrorMsg(ctx, "error.e0211"))
		return
	}
	external := httpx.GetExternalUser(ctx)
	if external == nil {
		httpx.WriteJSON(ctx, httpx.JsonErrorMsg(ctx, "error.e0150"))
		return
	}

	item := services.ConversationService.Get(id)
	if item == nil {
		httpx.WriteJSON(ctx, httpx.JsonErrorMsg(ctx, "error.e0116"))
		return
	}
	if !services.ConversationService.IsCustomerConversationOwner(item, *external) {
		httpx.WriteJSON(ctx, httpx.JsonErrorMsg(ctx, "error.e0222"))
		return
	}

	conversation := builders.BuildConversationWithLocale(item, i18nx.Locale(ctx))
	if channel := services.ChannelService.Get(item.ChannelID); channel != nil {
		conversation.ChannelName = channel.Name
	}
	detail := response.ConversationDetailResponse{
		ConversationResponse: conversation,
		Participants:         builders.BuildParticipantResponses(id),
	}
	httpx.WriteJSON(ctx, detail)
}

func ConversationPostCreate_or_match(ctx *gin.Context) {
	channel := services.ChannelService.GetEnabledChannel(ctx)
	if channel == nil {
		httpx.WriteJSON(ctx, httpx.JsonErrorMsg(ctx, "error.e0211"))
		return
	}
	external := httpx.GetExternalUser(ctx)
	if external == nil {
		httpx.WriteJSON(ctx, httpx.JsonErrorMsg(ctx, "error.e0150"))
		return
	}

	item, err := services.ConversationService.Create(*external, channel.ID, channel.AIAgentID)
	if err != nil {
		httpx.WriteJSON(ctx, err)
		return
	}
	conversation := builders.BuildConversationWithLocale(item, i18nx.Locale(ctx))
	conversation.ChannelName = channel.Name
	httpx.WriteJSON(ctx, conversation)
}

func ConversationPostClose(ctx *gin.Context) {
	if services.ChannelService.GetEnabledChannel(ctx) == nil {
		httpx.WriteJSON(ctx, httpx.JsonErrorMsg(ctx, "error.e0211"))
		return
	}
	external := httpx.GetExternalUser(ctx)
	if external == nil {
		httpx.WriteJSON(ctx, httpx.JsonErrorMsg(ctx, "error.e0150"))
		return
	}

	req := request.CloseConversationRequest{}
	if err := params.ReadJSON(ctx, &req); err != nil {
		httpx.WriteJSON(ctx, err)
		return
	}
	if err := services.ConversationService.CloseCustomerConversation(req.ConversationID, *external); err != nil {
		httpx.WriteJSON(ctx, err)
		return
	}
	httpx.WriteJSON(ctx, nil)
}
