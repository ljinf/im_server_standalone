package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ljinf/im_server_standalone/pkg/jwt"
	"github.com/ljinf/im_server_standalone/pkg/log"
)

type Handler struct {
	logger *log.Logger
	jwt    *jwt.JWT
}

func NewHandler(
	logger *log.Logger,
	j *jwt.JWT,
) *Handler {
	return &Handler{
		logger: logger,
		jwt:    j,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
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
