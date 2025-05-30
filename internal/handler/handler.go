package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ljinf/im_server_standalone/pkg/jwt"
	"github.com/ljinf/im_server_standalone/pkg/log"
	"github.com/spf13/viper"
)

type Handler struct {
	logger *log.Logger
	jwt    *jwt.JWT
	conf   *viper.Viper
}

func NewHandler(
	logger *log.Logger,
	j *jwt.JWT,
	conf *viper.Viper,
) *Handler {
	return &Handler{
		logger: logger,
		jwt:    j,
		conf:   conf,
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
	PageNum  int64 `json:"page_num"`
	PageSize int64 `json:"page_size"`
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
