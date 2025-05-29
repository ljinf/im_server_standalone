package service

import (
	"context"
	"fmt"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/repository"
	"github.com/ljinf/im_server_standalone/pkg/contants"
	"github.com/ljinf/im_server_standalone/pkg/util"
	"time"
)

type CommunityDomainService interface {
	//时刻
	AddMoment(ctx context.Context, req *v1.AddMomentReq) error
	EditMoment(ctx context.Context, req *v1.AddMomentReq) error
	GetMomentList(ctx context.Context, req *v1.GetMomentListReq) []model.CommunityMomentResp

	//点赞
	AddMomentLike(ctx context.Context, req *v1.AddMomentLikeReq) error

	//评论
	AddMomentComment(ctx context.Context, req *v1.AddMomentCommentReq) error
	AddMomentCommentLike(ctx context.Context, req *v1.AddMomentCommentLikeReq) error
	GetMomentCommentList(ctx context.Context, req *v1.MomentCommentListReq) []v1.MomentCommentListResp
}

type communityDomainService struct {
	*Service
	repo repository.CommunityRepository
}

func (c *communityDomainService) AddMomentCommentLike(ctx context.Context, req *v1.AddMomentCommentLikeReq) error {
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
		//c.logger.Error("CreateCommentLikeErr", "err", err.Error())
		//return errcode.ErrServer
		return err
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
		//log.Error(ctx, "CreateMomentCommentCountErr", "err", err.Error())
	}

	return nil
}

func (c *communityDomainService) AddMoment(ctx context.Context, req *v1.AddMomentReq) error {
	momentId, err := c.sid.GenUint64()
	if err != nil {
		//log.Error(ctx, "SidGenUint64Err", "err", err.Error())
		//return errcode.ErrServer
		return err
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
		//log.Error(ctx, "CreateCommunityMomentErr", "err", err.Error())
		//return errcode.ErrServer
		return err
	}

	//todo 审核内容

	return nil
}

func (c *communityDomainService) EditMoment(ctx context.Context, req *v1.AddMomentReq) error {
	moment := model.CommunityMoment{
		MomentId: req.MomentId,
		UserId:   req.UserId,
		Public:   req.Public,
		Status:   req.Status,
	}

	err := c.repo.UpdateCommunityMoment(ctx, &moment)
	if err != nil {
		//log.Error(ctx, "UpdateCommunityMoment", "err", err.Error())
		//return errcode.ErrServer
		return err
	}
	return nil
}

func (c *communityDomainService) GetMomentList(ctx context.Context, req *v1.GetMomentListReq) []model.CommunityMomentResp {

	list, err := c.repo.SelectCommunityMomentList(ctx, req)
	if err != nil {
		//log.Error(ctx, "SelectCommunityMomentListErr", "err", err.Error())
	}

	if len(list) > 0 {
		likes, err := c.repo.SelectUserMomentLikeList(ctx, req.UserId, list[len(list)-1].CreatedAt)
		if err != nil {
			//log.Error(ctx, "SelectUserMomentLikeListErr", "err", err.Error())
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

	return list
}

func (c *communityDomainService) AddMomentLike(ctx context.Context, req *v1.AddMomentLikeReq) error {
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
		//log.Error(ctx, "CreateMomentLikeErr", "err", err.Error())
		//return errcode.ErrServer
		return err
	}

	if err := c.repo.CreateMomentCount(ctx, &likeCount); err != nil {
		//log.Error(ctx, "CreateMomentCountErr", "err", err.Error())
	}

	return nil
}

func (c *communityDomainService) AddMomentComment(ctx context.Context, req *v1.AddMomentCommentReq) error {

	var (
		count int = 1
	)

	commentId, err := c.sid.GenUint64()
	if err != nil {
		//log.Error(ctx, "SidGenUint64Err", "err", err.Error())
		//return errcode.ErrServer
		return err
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
		//log.Error(ctx, "CreateMomentCommentErr", "err", err.Error())
		//return errcode.ErrServer
		return err
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
		//log.Error(ctx, "CreateMomentCountErr", "err", err.Error())
	}

	return nil
}

func (c *communityDomainService) GetMomentCommentList(ctx context.Context, req *v1.MomentCommentListReq) []v1.MomentCommentListResp {

	comments, err := c.repo.SelectMomentComment(ctx, req)
	if err != nil {
		//log.Error(ctx, "SelectMomentCommentErr", "err", err.Error())
		return []v1.MomentCommentListResp{}
	}

	resp := make([]v1.MomentCommentListResp, 0, len(comments))
	if err = util.CopyProperties(&resp, comments); err != nil {
		//log.Error(ctx, errcode.ErrCoverData.WithCause(err).Error())
	}

	if len(resp) > 0 {
		//点赞评论的回显
		likes, err := c.repo.SelectUserCommentLikeList(ctx, req.UserId, resp[0].CreatedAt)
		if err != nil {
			//log.Error(ctx, "SelectUserCommentLikeList", "err", err.Error())
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

func NewCommunityDomainService(d *Service, rep repository.CommunityRepository) CommunityDomainService {
	return &communityDomainService{
		Service: d,
		repo:    rep,
	}
}
