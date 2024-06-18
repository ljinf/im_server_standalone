package service

import (
	"context"
	"errors"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/repository"
	"github.com/ljinf/im_server_standalone/pkg/contants"
	"go.uber.org/zap"
	"time"
)

type RelationshipService interface {
	AddApplyFriendship(ctx context.Context, req *v1.ApplyFriendshipRequest) error
	GetApplyFriendshipList(ctx context.Context, userId int64, page int, pageSize int) (interface{}, error)
	UpdateApplyFriendshipInfo(ctx context.Context, req *v1.ApplyFriendshipRequest) error
	DelApplyFriendshipInfo(ctx context.Context, req *v1.ApplyFriendshipRequest) error

	GetRelationshipList(ctx context.Context, userId int64, page int, pageSize int) (interface{}, error)
	GetRelationship(ctx context.Context, req *v1.RelationshipRequest) (*model.RelationshipList, error)
	UpdateRelationshipList(ctx context.Context, req *v1.RelationshipRequest) error
	DelRelationship(ctx context.Context, req *v1.RelationshipRequest) error
}

type relationshipService struct {
	*Service
	repo repository.RelationshipRepository
}

func NewRelationshipService(s *Service, repo repository.RelationshipRepository) RelationshipService {
	return &relationshipService{
		Service: s,
		repo:    repo,
	}
}

func (r *relationshipService) AddApplyFriendship(ctx context.Context, req *v1.ApplyFriendshipRequest) error {
	now := time.Now()
	apply := model.ApplyFriendshipList{
		UserId:      req.UserId,
		TargetId:    req.TargetId,
		Remark:      req.Remark,
		Description: req.Description,
		Status:      req.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := r.repo.CreateApplyFriendship(ctx, &apply); err != nil {
		r.logger.Error(err.Error(), zap.Any("req", apply))
		return v1.ErrAddApplyFriendshipFailed
	}

	return nil
}

func (r *relationshipService) GetApplyFriendshipList(ctx context.Context, userId int64, page int, pageSize int) (interface{}, error) {
	list, total, err := r.repo.SelectRelationshipList(ctx, userId, contants.RelationshipTypeFriend, page, pageSize)
	if err != nil {
		r.logger.Error(err.Error(), zap.Any("userId", userId))
		return nil, v1.ErrInternalServerError
	}
	resp := map[string]interface{}{
		"rows":  list,
		"total": total,
	}
	return resp, nil
}

func (r *relationshipService) UpdateApplyFriendshipInfo(ctx context.Context, req *v1.ApplyFriendshipRequest) error {
	err := r.tm.Transaction(ctx, func(ctx context.Context) error {
		apply := model.ApplyFriendshipList{
			UserId:   req.TargetId,
			TargetId: req.UserId,
			Status:   req.Status,
		}

		if err := r.repo.UpdateApplyFriendship(ctx, &apply); err != nil {
			r.logger.Error(err.Error(), zap.Any("req", apply))
			return v1.ErrInternalServerError
		}

		now := time.Now()
		if apply.Status == contants.ApplyFriendshipStatusApproved {
			applyInfo, err := r.repo.SelectApplyOne(ctx, req.TargetId, req.UserId)
			if err != nil {
				r.logger.Error(err.Error(), zap.Any("userId", req.TargetId), zap.Any("targetId", req.UserId))
				if errors.Is(err, v1.ErrNotFound) {
					return v1.ErrBadRequest
				}
			}
			friendA := model.RelationshipList{
				UserId:           req.UserId,
				TargetId:         req.TargetId,
				Remark:           req.Remark,
				RelationshipType: contants.RelationshipTypeFriend,
				Status:           contants.RelationshipStatusNormal,
				CreatedAt:        now,
				UpdatedAt:        now,
			}

			friendB := model.RelationshipList{
				UserId:           req.TargetId,
				TargetId:         req.UserId,
				Remark:           applyInfo.Remark,
				RelationshipType: contants.RelationshipTypeFriend,
				Status:           contants.RelationshipStatusNormal,
				CreatedAt:        now,
				UpdatedAt:        now,
			}
			err = r.repo.CreateRelationship(ctx, friendA, friendB)
			if err != nil {
				r.logger.Error(err.Error(), zap.Any("friendship", [2]model.RelationshipList{friendA, friendB}))
				return v1.ErrCreateRelationshipFailed
			}
			return nil
		}
		return nil
	})
	return err
}

func (r *relationshipService) DelApplyFriendshipInfo(ctx context.Context, req *v1.ApplyFriendshipRequest) error {
	if err := r.repo.DelApplyFriendship(ctx, req.UserId, req.TargetId); err != nil {
		r.logger.Error(err.Error(), zap.Any("req", req))
		return v1.ErrInternalServerError
	}
	return nil
}

// 查询列表
func (r *relationshipService) GetRelationshipList(ctx context.Context, userId int64, page int, pageSize int) (interface{}, error) {
	list, total, err := r.repo.SelectRelationshipList(ctx, userId, contants.RelationshipTypeFriend, page, pageSize)
	if err != nil {
		r.logger.Error(err.Error(), zap.Any("userID", userId))
		return nil, err
	}
	resp := map[string]interface{}{
		"rows":  list,
		"total": total,
	}
	return resp, nil
}

// 查询一个
func (r *relationshipService) GetRelationship(ctx context.Context, req *v1.RelationshipRequest) (*model.RelationshipList, error) {
	info, err := r.repo.SelectRelationshipOne(ctx, req.UserId, req.TargetId, req.RelationshipType)
	if err != nil {
		r.logger.Error(err.Error(), zap.Any("param", req))
		return nil, err
	}
	return info, nil
}

func (r *relationshipService) UpdateRelationshipList(ctx context.Context, req *v1.RelationshipRequest) error {
	info := model.RelationshipList{
		UserId:   req.UserId,
		TargetId: req.TargetId,
		Remark:   req.Remark,
		Status:   req.Status,
	}
	if err := r.repo.UpdateRelationship(ctx, &info); err != nil {
		r.logger.Error(err.Error(), zap.Any("req", info))
		return v1.ErrInternalServerError
	}
	return nil
}

func (r *relationshipService) DelRelationship(ctx context.Context, req *v1.RelationshipRequest) error {
	return r.repo.DelRelationship(ctx, req.UserId, req.TargetId, req.RelationshipType)
}
