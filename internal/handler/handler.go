package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ljinf/im_server_standalone/pkg/jwt"
	"github.com/ljinf/im_server_standalone/pkg/log"
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
	v, exists := ctx.Get("claims")
	if !exists {
		return 0
	}
	userInfo := v.(*jwt.MyCustomClaims)
	return userInfo.UserId
}

type PageInfo struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

func GetPageInfo(ctx *gin.Context) *PageInfo {
	var pageInfo PageInfo
	_ = ctx.ShouldBind(&pageInfo)

	if pageInfo.PageNum == 0 {
		pageInfo.PageNum = 1
	}

	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 30
	}

	return &pageInfo
}
