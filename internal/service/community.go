package service

import (
	"context"
	"fmt"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/repository"
	"github.com/ljinf/im_server_standalone/pkg/contants"
	"github.com/ljinf/im_server_standalone/pkg/util"
	"go.uber.org/zap"
	"time"
)

type CommunityService interface {
	//时刻
	AddMoment(ctx context.Context, req *v1.AddMomentReq) error
	EditMoment(ctx context.Context, req *v1.AddMomentReq) error
	GetMomentList(ctx context.Context, req *v1.MomentListReq) []v1.MomentListResp

	//点赞
	AddMomentLike(ctx context.Context, req *v1.AddMomentLikeReq) error

	//评论
	AddMomentComment(ctx context.Context, req *v1.AddMomentCommentReq) (*v1.MomentCommentListResp, error)
	AddMomentCommentLike(ctx context.Context, req *v1.AddMomentCommentLikeReq) error
	GetMomentCommentList(ctx context.Context, req *v1.MomentCommentListReq) []v1.MomentCommentListResp
}

type communityService struct {
	*Service
	repo repository.CommunityRepository
}

func (c *communityService) AddMomentCommentLike(ctx context.Context, req *v1.AddMomentCommentLikeReq) error {
	var (
		count       int
		cancelCount int
	)

	if err := c.repo.CreateCommentLike(ctx, &model.MomentCommentLike{
		UserId:    req.UserId,
		CommentId: req.CommentId,
		Status:    req.Status,
		CreatedAt: time.Now().Unix(),
	}); err != nil {
		c.logger.Error(err.Error(), zap.Any("CreateCommentLike", req))
		return v1.ErrInternalServerError
	}

	switch req.Status {
	case contants.MomentLikeNormal:
		count = 1
		break
	case contants.MomentLikeCancel:
		cancelCount = 1
		break
	}

	countInfo := model.MomentCommentCount{
		CommentId:       req.CommentId,
		LikeCancelCount: cancelCount,
		LikeCount:       count,
	}
	if err := c.repo.CreateMomentCommentCount(ctx, &countInfo); err != nil {
		c.logger.Error(err.Error(), zap.Any("CreateMomentCommentCountErr", countInfo))
	}

	return nil
}

func (c *communityService) AddMoment(ctx context.Context, req *v1.AddMomentReq) error {
	momentId, err := c.sid.GenUint64()
	if err != nil {
		c.logger.Error(err.Error(), zap.Any("GenUint64", ""))
		return v1.ErrInternalServerError
	}
	moment := model.CommunityMoment{
		UserId:         req.UserId,
		MomentId:       fmt.Sprintf("%v", momentId),
		Content:        req.Content,
		Attachment:     req.Attachment,
		AttachmentType: req.AttachmentType,
		Public:         req.Public,
		Status:         contants.MomentStatusPending,
		CreatedAt:      time.Now().Unix(),
	}

	err = c.repo.CreateCommunityMoment(ctx, &moment)
	if err != nil {
		c.logger.Error(err.Error(), zap.Any("CreateCommunityMoment", ""))
		return v1.ErrInternalServerError
	}

	//todo 审核内容

	return nil
}

func (c *communityService) EditMoment(ctx context.Context, req *v1.AddMomentReq) error {
	moment := model.CommunityMoment{
		MomentId: req.MomentId,
		UserId:   req.UserId,
		Public:   req.Public,
		Status:   req.Status,
	}

	err := c.repo.UpdateCommunityMoment(ctx, &moment)
	if err != nil {
		c.logger.Error(err.Error(), zap.Any("UpdateCommunityMoment", ""))
		return v1.ErrInternalServerError
	}
	return nil
}

func (c *communityService) GetMomentList(ctx context.Context, req *v1.MomentListReq) []v1.MomentListResp {

	list, err := c.repo.SelectCommunityMomentList(ctx, req)
	if err != nil {
		c.logger.Error(err.Error(), zap.Any("SelectCommunityMomentList", ""))
	}

	if len(list) > 0 {
		likes, err := c.repo.SelectUserMomentLikeList(ctx, req.UserId, list[len(list)-1].CreatedAt)
		if err != nil {
			c.logger.Error(err.Error(), zap.Any("SelectUserMomentLikeList", ""))
		}
		for index := range list {
			for _, v := range likes {
				if list[index].MomentId == v.MomentId {
					list[index].LikeStatus = 1
					break
				}
			}
		}
	}

	resp := make([]v1.MomentListResp, 0, len(list))
	if err = util.CopyProperties(&resp, list); err != nil {
		c.logger.Error(err.Error(), zap.Any("CopyProperties", ""))
	}

	return resp
}

func (c *communityService) AddMomentLike(ctx context.Context, req *v1.AddMomentLikeReq) error {
	var (
		count       int = 0
		cancelCount int = 0
	)
	momentLike := model.MomentLike{
		MomentId:  req.MomentId,
		UserId:    req.UserId,
		Status:    req.Status,
		CreatedAt: time.Now().Unix(),
	}

	switch req.Status {
	case contants.MomentLikeNormal:
		count = 1
		break
	case contants.MomentLikeCancel:
		cancelCount = 1
		break
	}
	likeCount := model.MomentCount{
		MomentId:        req.MomentId,
		LikeCount:       count,
		LikeCancelCount: cancelCount,
	}

	//这里不用事务，即使计数错误也不影响时刻的添加
	if err := c.repo.CreateMomentLike(ctx, &momentLike); err != nil {
		c.logger.Error(err.Error(), zap.Any("CreateMomentLike", ""))
		return v1.ErrInternalServerError
	}

	if err := c.repo.CreateMomentCount(ctx, &likeCount); err != nil {
		c.logger.Error(err.Error(), zap.Any("CreateMomentCount", ""))
	}

	return nil
}

func (c *communityService) AddMomentComment(ctx context.Context, req *v1.AddMomentCommentReq) (*v1.MomentCommentListResp, error) {

	var (
		count int = 1
	)

	commentId, err := c.sid.GenUint64()
	if err != nil {
		c.logger.Error(err.Error(), zap.Any("AddMomentComment", ""))
		return nil, v1.ErrInternalServerError
	}

	commentInfo := &model.MomentComment{
		CommentId:      fmt.Sprintf("%v", commentId),
		ParentId:       req.ParentId,
		MomentId:       req.MomentId,
		UserId:         req.UserId,
		ReplyId:        req.ReplyId,
		ReplyCommentId: req.ReplyCommentId,
		Content:        req.Content,
		Status:         contants.MomentCommentStatusNormal,
		CreatedAt:      time.Now().Unix(),
	}

	if err = c.repo.CreateMomentComment(ctx, commentInfo); err != nil {
		c.logger.Error(err.Error(), zap.Any("AddMomentComment", ""))
		return nil, v1.ErrInternalServerError
	}

	//计数统计
	if err = c.tm.Transaction(ctx, func(ctx context.Context) error {
		//评论数
		if err = c.repo.CreateMomentCount(ctx, &model.MomentCount{
			MomentId:     req.MomentId,
			CommentCount: count,
		}); err != nil {
			return err
		}

		//回复数
		if commentInfo.ParentId != "0" {
			return c.repo.CreateMomentCommentCount(ctx, &model.MomentCommentCount{
				CommentId:    commentInfo.ParentId,
				CommentCount: count,
			})
		}

		return nil
	}); err != nil {
		c.logger.Error(err.Error(), zap.Any("AddMomentComment", ""))
	}

	return &v1.MomentCommentListResp{
		CommentId:       commentInfo.CommentId,
		ParentId:        commentInfo.ParentId,
		MomentId:        commentInfo.MomentId,
		UserId:          commentInfo.UserId,
		ReplyId:         commentInfo.ReplyId,
		Content:         commentInfo.Content,
		LikeCount:       0,
		LikeStatus:      0,
		LikeCancelCount: 0,
		CommentCount:    0,
		Status:          commentInfo.Status,
		CreatedAt:       commentInfo.CreatedAt,
	}, nil
}

func (c *communityService) GetMomentCommentList(ctx context.Context, req *v1.MomentCommentListReq) []v1.MomentCommentListResp {

	comments, err := c.repo.SelectMomentComment(ctx, req)
	if err != nil {
		c.logger.Error(err.Error(), zap.Any("GetMomentCommentList", req))
		return []v1.MomentCommentListResp{}
	}

	resp := make([]v1.MomentCommentListResp, 0, len(comments))
	if err = util.CopyProperties(&resp, comments); err != nil {
		c.logger.Error(err.Error(), zap.Any("GetMomentCommentList", req))
	}

	if len(resp) > 0 {
		//点赞评论的回显
		likes, err := c.repo.SelectUserCommentLikeList(ctx, req.UserId, resp[0].CreatedAt)
		if err != nil {
			c.logger.Error(err.Error(), zap.Any("SelectUserCommentLikeList",
				fmt.Sprintf("%v %v", req.UserId, resp[0].CreatedAt)))
		}
		for index := range likes {
			for _, v := range resp {
				if resp[index].CommentId == v.CommentId {
					resp[index].LikeStatus = 1
					break
				}
			}
		}
	}

	return resp
}

func NewCommunityService(d *Service, rep repository.CommunityRepository) CommunityService {
	return &communityService{
		Service: d,
		repo:    rep,
	}
}
