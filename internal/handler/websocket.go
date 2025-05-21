package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/service"
	"go.uber.org/zap"
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

	if err := h.authorization(ctx); err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, err, nil)
		return
	}

	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
	userId := GetUserIdFromCtx(ctx)
	h.srv.InitConn(userId, conn)
}

func (h *webSocketHandler) authorization(ctx *gin.Context) error {
	token := ctx.Query("token")
	claims, err := h.jwt.ParseToken(token)
	if err != nil {
		h.logger.WithContext(ctx).Error("token error", zap.Any("data", map[string]interface{}{
			"url":    ctx.Request.URL,
			"params": ctx.Params,
		}), zap.Error(err))
		return v1.ErrUnauthorized
	}
	ctx.Set("claims", claims)
	return nil
}
