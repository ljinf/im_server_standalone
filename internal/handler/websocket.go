package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ljinf/im_server_standalone/internal/service"
	"net/http"
)

var (
	wsUpgrader = websocket.Upgrader{
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WebSocketHandler interface {
	AcceptConn(ctx *gin.Context)
}

type webSocketHandler struct {
	*Handler
	srv service.WebsocketService
}

func NewWebSocketHandler(h *Handler, s service.WebsocketService) WebSocketHandler {
	return &webSocketHandler{
		Handler: h,
		srv:     s,
	}
}

func (h *webSocketHandler) AcceptConn(ctx *gin.Context) {
	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
	userId := GetUserIdFromCtx(ctx)
	h.srv.InitConn(userId, conn)
}
