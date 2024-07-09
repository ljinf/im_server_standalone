package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ljinf/im_server_standalone/internal/service"
)

type ChatHandler struct {
	*Handler
	srv service.ChatService
}

func NewChatHandler(h *Handler, srv service.ChatService) *ChatHandler {
	return &ChatHandler{
		Handler: h,
		srv:     srv,
	}
}

// 发送信息
func (h *ChatHandler) SendChatMessage(ctx *gin.Context) {

}

// 会话列表
func (h *ChatHandler) GetUserConversationList(ctx *gin.Context) {

}

// 消息列表
func (h *ChatHandler) GetUserMsgList(ctx *gin.Context) {

}
