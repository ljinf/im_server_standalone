package handler

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type CommunityHandler interface {
	//时刻
	AddMoment(ctx *gin.Context)
	EditMoment(ctx *gin.Context)
	GetMomentList(ctx *gin.Context)

	//点赞
	AddMomentLike(ctx *gin.Context)

	//评论
	AddMomentComment(ctx *gin.Context)
	//评论被点赞
	LikeMomentComment(ctx *gin.Context)
	//评论列表
	GetMomentCommentList(ctx *gin.Context)
}

type communityHandler struct {
	*Handler
	srv service.CommunityService
}

func (h *communityHandler) AddMoment(ctx *gin.Context) {
	var param v1.AddMomentReq
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error(), zap.Any("AddMoment", ""))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if len(param.Attachment) > 0 && param.AttachmentType == 0 {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	param.UserId = GetUserIdFromCtx(ctx)

	if err := h.srv.AddMoment(ctx, &param); err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *communityHandler) EditMoment(ctx *gin.Context) {
	var param v1.AddMomentReq
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error(), zap.Any("AddMoment", ""))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if len(param.MomentId) == 0 {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	param.UserId = GetUserIdFromCtx(ctx)

	if err := h.srv.EditMoment(ctx, &param); err != nil {
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *communityHandler) GetMomentList(ctx *gin.Context) {
	var param v1.MomentListReq
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error(), zap.Any("AddMoment", ""))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	param.UserId = GetUserIdFromCtx(ctx)
	pageInfo := GetPageInfo(ctx)
	param.PageNum = int(pageInfo.PageNum)
	param.PageSize = int(pageInfo.PageSize)

	momentList := h.srv.GetMomentList(ctx, &param)
	v1.HandleSuccess(ctx, momentList)
}

func (h *communityHandler) AddMomentLike(ctx *gin.Context) {
	var param v1.AddMomentLikeReq
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error(), zap.Any("AddMoment", ""))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	param.UserId = GetUserIdFromCtx(ctx)

	if err := h.srv.AddMomentLike(ctx, &param); err != nil {
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *communityHandler) AddMomentComment(ctx *gin.Context) {
	var param v1.AddMomentCommentReq
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error(), zap.Any("AddMoment", ""))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	param.UserId = GetUserIdFromCtx(ctx)

	if err := h.srv.AddMomentComment(ctx, &param); err != nil {
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *communityHandler) GetMomentCommentList(ctx *gin.Context) {
	var param v1.MomentCommentListReq
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error(), zap.Any("AddMoment", ""))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	param.UserId = GetUserIdFromCtx(ctx)
	pageInfo := GetPageInfo(ctx)
	param.PageNum = int(pageInfo.PageNum)
	param.PageSize = int(pageInfo.PageSize)

	list := h.srv.GetMomentCommentList(ctx, &param)

	v1.HandleSuccess(ctx, list)
}

func (h *communityHandler) LikeMomentComment(ctx *gin.Context) {
	var param v1.AddMomentCommentLikeReq
	if err := ctx.ShouldBind(&param); err != nil {
		h.logger.Error(err.Error(), zap.Any("AddMoment", ""))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	param.UserId = GetUserIdFromCtx(ctx)

	if err := h.srv.AddMomentCommentLike(ctx, &param); err != nil {
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func NewCommunityHandler(h *Handler, srv service.CommunityService) CommunityHandler {
	return &communityHandler{
		Handler: h,
		srv:     srv,
	}
}
