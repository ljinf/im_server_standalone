package handler

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type ChatHandler struct {
	*Handler
	srv       service.ChatService
	socketSrv service.WebsocketService
}

func NewChatHandler(h *Handler, srv service.ChatService, socket service.WebsocketService) *ChatHandler {
	return &ChatHandler{
		Handler:   h,
		srv:       srv,
		socketSrv: socket,
	}
}

// 发送信息
func (h *ChatHandler) SendChatMessage(ctx *gin.Context) {
	var params v1.SendMsgReq
	if err := ctx.ShouldBind(&params); err != nil {
		h.logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}

	msgResp, err := h.srv.CreateMsg(ctx, &params)
	if err != nil {
		h.logger.Error(err.Error(), zap.Any("param", params))
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	// 转发给target
	h.socketSrv.SyncPushMsg(msgResp, params.TargetId)

	v1.HandleSuccess(ctx, msgResp)
}

// 会话列表
func (h *ChatHandler) GetUserConversationList(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}

	conversationList, err := h.srv.GetUserConversationList(ctx, userId)
	if err != nil {
		h.logger.Error(err.Error(), zap.Any("userId", userId))
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, conversationList)
}

// 消息列表
func (h *ChatHandler) GetUserMsgList(ctx *gin.Context) {
	var params v1.HistoryMsgListReq
	if err := ctx.ShouldBind(&params); err != nil {
		h.logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}

	msgList, err := h.srv.GetMsgList(ctx, params.UserId, params.ConversationId, params.Seq, params.PageNum, params.PageSize)
	if err != nil {
		h.logger.Error(err.Error(), zap.Any("params", params))
		v1.HandleSuccess(ctx, []model.MsgList{})
		return
	}
	v1.HandleSuccess(ctx, msgList)
}

// 上报已读
func (h *ChatHandler) ReportReadMsgSeq(ctx *gin.Context) {
	var params v1.ReportReadReq
	if err := ctx.ShouldBind(&params); err != nil {
		h.logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}

	if err := h.srv.ReportReadMsgSeq(ctx, &params); err != nil {
		h.logger.Error(err.Error(), zap.Any("params", params))
	}
	v1.HandleSuccess(ctx, nil)
}
