package handler

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/service"
	"net/http"
)

type RelationshipHandler struct {
	*Handler
	srv service.RelationshipService
}

func NewRelationshipHandler(h *Handler, srv service.RelationshipService) *RelationshipHandler {
	return &RelationshipHandler{
		Handler: h,
		srv:     srv,
	}
}

// 申请好友
func (h *RelationshipHandler) AddApplyFriendship(ctx *gin.Context) {
	var param v1.ApplyFriendshipRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}

	if err := h.srv.AddApplyFriendship(ctx, &param); err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// 获取申请列表
func (h *RelationshipHandler) GetApplyFriendshipList(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}

	pageInfo := GetPageInfo(ctx)
	list, err := h.srv.GetApplyFriendshipList(ctx, userId, pageInfo.PageNum, pageInfo.PageSize)
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, list)
}

// 修改申请信息
func (h *RelationshipHandler) UpdateApplyFriendshipInfo(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}

	var param v1.ApplyFriendshipRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}

	param.UserId = userId
	if err := h.srv.UpdateApplyFriendshipInfo(ctx, &param); err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// 删除申请
func (h *RelationshipHandler) DelApplyFriendshipInfo(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}

	var param v1.ApplyFriendshipRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}
	param.UserId = userId
	if err := h.srv.DelApplyFriendshipInfo(ctx, &param); err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// 获取关系列表
func (h *RelationshipHandler) GetRelationshipList(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}

	pageInfo := GetPageInfo(ctx)
	list, err := h.srv.GetRelationshipList(ctx, userId, pageInfo.PageNum, pageInfo.PageSize)
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, list)
}

// 关注
func (h *RelationshipHandler) AddRelationshipFollow(ctx *gin.Context) {

}

// 修改关系
func (h *RelationshipHandler) UpdateRelationship(ctx *gin.Context) {

}

// 删除关系
func (h *RelationshipHandler) DelRelationship(ctx *gin.Context) {

}
