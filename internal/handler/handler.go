package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ljinf/im_server_standalone/pkg/log"
	"time"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(
	logger *log.Logger,
) *Handler {
	return &Handler{
		logger: logger,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) int64 {
	//v, exists := ctx.Get("claims")
	//if !exists {
	//	return ""
	//}
	return time.Now().Unix()
}
