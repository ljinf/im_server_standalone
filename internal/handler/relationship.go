package handler

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/service"
	"github.com/ljinf/im_server_standalone/pkg/contants"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type RelationshipHandler struct {
	*Handler
	srv   service.RelationshipService
	imSrv service.ChatService
}

func NewRelationshipHandler(h *Handler, srv service.RelationshipService, imsvr service.ChatService) *RelationshipHandler {
	return &RelationshipHandler{
		Handler: h,
		srv:     srv,
		imSrv:   imsvr,
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

	if param.Status == contants.ApplyFriendshipStatusApproved {
		msgReq := &v1.SendMsgReq{
			UserId:      userId,
			TargetId:    param.TargetId,
			Content:     contants.ChatSayHello,
			ContentType: contants.MsgContentTypeTxt,
			SendTime:    time.Now().Unix(),
			CreatedAt:   time.Now().Unix(),
		}
		_, err := h.imSrv.CreateMsg(ctx, msgReq)
		if err != nil {
			h.logger.Error(err.Error())
		}
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

	var param v1.RelationshipListReq
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}

	if param.PageNum == 0 {
		param.PageNum = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 30
	}

	list, err := h.srv.GetRelationshipList(ctx, userId, param.RelationshipType, param.PageNum, param.PageSize)
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, list)
}

// 关注
func (h *RelationshipHandler) AddRelationshipFollow(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}

	var param v1.RelationshipRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}
	param.UserId = userId
	param.RelationshipType = contants.RelationshipTypeFollow

	if err := h.srv.AddRelationshipFollow(ctx, &param); err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// 修改关系
func (h *RelationshipHandler) UpdateRelationship(ctx *gin.Context) {

	userId := GetUserIdFromCtx(ctx)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}

	var param v1.RelationshipRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}
	param.UserId = userId

	if err := h.srv.UpdateRelationship(ctx, &param); err != nil {
		h.logger.Error(err.Error(), zap.Any("param", param))
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// 删除关系
func (h *RelationshipHandler) DelRelationship(ctx *gin.Context) {

	userId := GetUserIdFromCtx(ctx)
	if userId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}

	var param v1.RelationshipRequest
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}

	param.UserId = userId
	if err := h.srv.DelRelationship(ctx, &param); err != nil {
		h.logger.Error(err.Error(), zap.Any("param", param))
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}
