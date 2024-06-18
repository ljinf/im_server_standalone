package repository

import (
	"context"
	"errors"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type RelationshipRepository interface {
	// 好友申请相关
	CreateApplyFriendship(ctx context.Context, info *model.ApplyFriendshipList) error
	UpdateApplyFriendship(ctx context.Context, info *model.ApplyFriendshipList) error
	SelectApplyFriendshipList(ctx context.Context, userId int64, page, pageSize int) ([]model.ApplyFriendshipList, error)
	SelectApplyOne(ctx context.Context, userId, targetId int64) (*model.ApplyFriendshipList, error)
	DelApplyFriendship(ctx context.Context, userId, targetId int64) error

	// 关系
	CreateRelationship(ctx context.Context, list ...model.RelationshipList) error
	SelectRelationshipList(ctx context.Context, userId int64, relationshipType, page, pageSize int) ([]model.RelationshipList, int, error)
	SelectRelationshipOne(ctx context.Context, userId, targetId int64, relationshipType int) (*model.RelationshipList, error)
	UpdateRelationship(ctx context.Context, info *model.RelationshipList) error
	DelRelationship(ctx context.Context, userId, targetId int64, relationshipType int) error
}

type relationshipRepository struct {
	*Repository
}

func NewRelationshipRepository(repo *Repository) RelationshipRepository {
	return &relationshipRepository{
		Repository: repo,
	}
}

// 申请相关
func (r *relationshipRepository) CreateApplyFriendship(ctx context.Context, info *model.ApplyFriendshipList) error {
	return r.DB(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "target_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"remark": info.Remark, "description": info.Description,
			"status": info.Status, "updated_at": time.Now()}),
	}).Create(info).Error
}

func (r *relationshipRepository) UpdateApplyFriendship(ctx context.Context, info *model.ApplyFriendshipList) error {
	return r.DB(ctx).Where("user_id=? and target_id=?", info.UserId, info.TargetId).Updates(
		map[string]interface{}{"status": info.Status, "updated_at": time.Now()}).Error
}

func (r *relationshipRepository) SelectApplyOne(ctx context.Context, userId, targetId int64) (*model.ApplyFriendshipList, error) {
	var info model.ApplyFriendshipList
	if err := r.DB(ctx).Where("user_id=? and target_id=?", userId, targetId).Find(&info).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}
	return &info, nil
}

func (r *relationshipRepository) SelectApplyFriendshipList(ctx context.Context, userId int64, page, pageSize int) ([]model.ApplyFriendshipList, error) {
	var list []model.ApplyFriendshipList
	if err := r.DB(ctx).Where("user_id=?", userId).Order("updated_at desc").
		Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *relationshipRepository) DelApplyFriendship(ctx context.Context, userId, targetId int64) error {
	return r.DB(ctx).Where("user_id=? and target_id=?", userId, targetId).Delete(&model.ApplyFriendshipList{}).Error
}

// 关系相关
func (r *relationshipRepository) CreateRelationship(ctx context.Context, list ...model.RelationshipList) error {
	if len(list) > 0 {
		result := r.DB(ctx).Create(&list)
		if result.Error != nil {
			return result.Error
		}
		if int(result.RowsAffected) != len(list) {
			return errors.New("create relationship failed")
		}
	}
	return nil
}

func (r *relationshipRepository) SelectRelationshipList(ctx context.Context, userId int64, relationshipType, page, pageSize int) ([]model.RelationshipList, int, error) {

	conds := []string{"user_id=?", "status=?", "relationship_type=?"}
	values := []interface{}{userId, 1, relationshipType}
	lists, total, err := r.doSelectRelationship(ctx, conds, values, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return lists, total, nil
}

func (r *relationshipRepository) SelectRelationshipOne(ctx context.Context, userId, targetId int64, relationshipType int) (*model.RelationshipList, error) {
	conds := []string{"user_id=?", "target_id=?", "relationship_type=?"}
	values := []interface{}{userId, targetId, relationshipType}
	lists, _, err := r.doSelectRelationship(ctx, conds, values, 1, 1)
	if err != nil {
		return nil, err
	}

	if len(lists) < 1 {
		return nil, v1.ErrNotFound
	}

	return &lists[0], nil
}

func (r *relationshipRepository) doSelectRelationship(ctx context.Context, conds []string, values []interface{}, page, pageSize int) ([]model.RelationshipList, int, error) {
	var list []model.RelationshipList

	if err := r.DB(ctx).Where(strings.Join(conds, " and "), values...).
		Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := r.DB(ctx).Model(&model.RelationshipList{}).Where(strings.Join(conds, " and "), values...).Count(&count).Error; err != nil {
		r.logger.Error(err.Error(), zap.Any("conds", conds), zap.Any("values", values))
	}
	return list, int(count), nil
}

func (r *relationshipRepository) UpdateRelationship(ctx context.Context, info *model.RelationshipList) error {
	return r.DB(ctx).Where("user_id=? and target_id=? and relationship_type=?",
		info.UserId, info.TargetId, info.RelationshipType).Updates(info).Error
}

func (r *relationshipRepository) DelRelationship(ctx context.Context, userId, targetId int64, relationshipType int) error {
	return r.DB(ctx).Where("user_id=? and target_id=? and relationship_type=?", userId, targetId, relationshipType).
		Delete(&model.RelationshipList{}).Error
}
