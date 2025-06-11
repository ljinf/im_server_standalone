package repository

import (
	"context"
	"errors"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/pkg/contants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type CommunityRepository interface {
	//用户时刻
	CreateCommunityMoment(ctx context.Context, req *model.CommunityMoment) error
	UpdateCommunityMoment(ctx context.Context, req *model.CommunityMoment) error
	SelectCommunityMomentList(ctx context.Context, params *v1.MomentListReq) ([]model.CommunityMomentResp, error)

	//点赞
	CreateMomentLike(ctx context.Context, req *model.MomentLike) error
	CreateCommentLike(ctx context.Context, req *model.MomentCommentLike) error
	//用户点赞记录
	SelectUserMomentLikeList(ctx context.Context, userId string, momentCreateAt int64) ([]model.MomentLike, error)
	SelectUserCommentLikeList(ctx context.Context, userId string, momentCreateAt int64) ([]model.MomentCommentLike, error)
	//计数
	CreateMomentCount(ctx context.Context, req *model.MomentCount) error
	SelectMomentCount(ctx context.Context, momentId string) (*model.MomentCount, error)

	//评论
	CreateMomentComment(ctx context.Context, req *model.MomentComment) error
	UpdateMomentComment(ctx context.Context, req *model.MomentComment) error
	CreateMomentCommentCount(ctx context.Context, req *model.MomentCommentCount) error //点赞计数
	SelectMomentComment(ctx context.Context, params *v1.MomentCommentListReq) ([]model.MomentCommentResp, error)
}

type communityRepository struct {
	*Repository
}

func (r *communityRepository) CreateCommunityMoment(ctx context.Context, req *model.CommunityMoment) error {
	return r.DB(ctx).Create(req).Error
}

func (r *communityRepository) UpdateCommunityMoment(ctx context.Context, req *model.CommunityMoment) error {
	return r.DB(ctx).Where("user_id=? and moment_id=?").Updates(req).Error
}

func (r *communityRepository) SelectCommunityMomentList(ctx context.Context, params *v1.MomentListReq) ([]model.CommunityMomentResp, error) {

	if len(params.WhereUser) > 0 {
		return r.selectUserMomentList(ctx, params)
	}

	return r.selectPublicMomentList(ctx, params)
}

// 公共时刻
func (r *communityRepository) selectPublicMomentList(ctx context.Context, params *v1.MomentListReq) ([]model.CommunityMomentResp, error) {

	var (
		conds  = []string{"ml.`public`=1", "ml.`status`=2"}
		values = []interface{}{}
	)

	if params.CreatedAt != 0 {
		switch params.Direct {
		case contants.DirectLessThan:
			conds = append(conds, "ml.`created_at` < ?")
			break
		case contants.DirectGreaterThan:
			conds = append(conds, "ml.`created_at` > ?")
			break
		}

		values = append(values, params.CreatedAt)
	}

	var list []model.CommunityMomentResp

	querySql := "SELECT ml.`user_id`,ml.`moment_id`,ml.`content`,ml.`attachment`,ml.`attachment_type`,ml.`public`,ml.`status`,ml.`created_at`,mc.`like_count`,mc.`like_cancel_count`,mc.`comment_count` " +
		"FROM `community_moment_list` ml LEFT JOIN `moment_count_list` mc ON mc.`moment_id`=ml.`moment_id` " +
		" WHERE " + strings.Join(conds, " and ") +
		" ORDER BY ml.`id` DESC LIMIT ?"
	err := r.DB(ctx).Raw(querySql, append(values, params.PageSize)...).Find(&list).Error

	return list, err
}

// 个人时刻
func (r *communityRepository) selectUserMomentList(ctx context.Context, params *v1.MomentListReq) ([]model.CommunityMomentResp, error) {

	var (
		conds  = []string{"ml.`user_id` = ?"}
		values = []interface{}{params.WhereUser}
	)

	if params.CreatedAt != 0 {
		switch params.Direct {
		case contants.DirectLessThan:
			conds = append(conds, "ml.`created_at` < ?")
			break
		case contants.DirectGreaterThan:
			conds = append(conds, "ml.`created_at` > ?")
			break
		}

		values = append(values, params.CreatedAt)
	}

	//如果是查看别人的时刻
	if params.UserId != params.WhereUser {
		conds = append(conds, "ml.`public` = 1 and ml.`status` = 2")
	}

	var list []model.CommunityMomentResp

	querySql := "SELECT ml.`user_id`,ml.`moment_id`,ml.`content`,ml.`attachment`,ml.`attachment_type`,ml.`public`,ml.`status`,ml.`created_at`,mc.`like_count`,mc.`like_cancel_count`,mc.`comment_count` " +
		"FROM `community_moment_list` ml LEFT JOIN `moment_count_list` mc ON mc.`moment_id`=ml.`moment_id` " +
		" WHERE " + strings.Join(conds, " and ") +
		" ORDER BY ml.`moment_id` DESC LIMIT ? OFFSET ?"
	err := r.DB(ctx).Raw(querySql, append(values, params.PageSize, (params.PageNum-1)*params.PageSize)...).Find(&list).Error

	return list, err
}

func (r *communityRepository) CreateMomentLike(ctx context.Context, req *model.MomentLike) error {
	return r.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"status", "created_at"}),
	}).Create(req).Error
}

func (r *communityRepository) CreateCommentLike(ctx context.Context, req *model.MomentCommentLike) error {
	return r.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"status", "created_at"}),
	}).Create(req).Error
}

func (r *communityRepository) CreateMomentComment(ctx context.Context, req *model.MomentComment) error {
	return r.DB(ctx).Create(req).Error
}

func (r *communityRepository) UpdateMomentComment(ctx context.Context, req *model.MomentComment) error {
	return r.DB(ctx).Updates(req).Error
}

func (r *communityRepository) SelectMomentComment(ctx context.Context, params *v1.MomentCommentListReq) ([]model.MomentCommentResp, error) {
	var (
		conds  = []string{"1=1"}
		values = []interface{}{}
		list   []model.MomentCommentResp
	)

	if len(params.MomentId) > 0 {
		conds = append(conds, "mc.`moment_id`=?")
		values = append(values, params.MomentId)
	}

	if len(params.ParentId) > 0 {
		conds = append(conds, "mc.`parent_id`=?")
		values = append(values, params.ParentId)
	} else {
		//一级评论
		conds = append(conds, "mc.`parent_id`=?")
		values = append(values, 0)
	}

	querySql := "SELECT mc.`id`,mc.`comment_id`,mc.`parent_id`,mc.`moment_id`,mc.`user_id`,mc.`reply_id`,mc.`content`,mc.`status`,mc.`created_at`," +
		"IFNULL(cc.`like_count`,0)`like_count`,IFNULL(cc.`like_cancel_count`,0)`like_cancel_count`,IFNULL(cc.`comment_count`,0)`comment_count` " +
		"FROM `moment_comment_list` mc LEFT JOIN `moment_comment_count_list` cc ON cc.`comment_id`=mc.`comment_id` " +
		"WHERE " + strings.Join(conds, " and ") + " LIMIT ? OFFSET ?"

	err := r.DB(ctx).Raw(querySql, append(values, params.PageSize, (params.PageNum-1)*params.PageSize)...).Find(&list).Error

	return list, err
}

// momentCreateAt  时刻创建时间
func (r *communityRepository) SelectUserMomentLikeList(ctx context.Context, userId string, momentCreateAt int64) ([]model.MomentLike, error) {
	var list []model.MomentLike
	err := r.DB(ctx).Where("user_id=? and created_at>=? and status=1", userId, momentCreateAt).Find(&list).Error
	return list, err
}

// createAt  评论创建时间
func (r *communityRepository) SelectUserCommentLikeList(ctx context.Context, userId string, createAt int64) ([]model.MomentCommentLike, error) {
	var list []model.MomentCommentLike
	err := r.DB(ctx).Where("user_id=? and created_at>=? and status=1", userId, createAt).Find(&list).Error
	return list, err
}

func (r *communityRepository) CreateMomentCount(ctx context.Context, req *model.MomentCount) error {
	return r.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"like_count":        gorm.Expr("like_count+VALUES(like_count)"),
			"like_cancel_count": gorm.Expr("like_cancel_count+VALUES(like_cancel_count)"),
			"comment_count":     gorm.Expr("comment_count+VALUES(comment_count)"),
		}),
	}).Create(req).Error
}

func (r *communityRepository) SelectMomentCount(ctx context.Context, momentId string) (*model.MomentCount, error) {
	var item model.MomentCount
	err := r.DB(ctx).Where("moment_id=?", momentId).First(&item).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return &item, nil
	}
	return &item, err
}

func (r *communityRepository) CreateMomentCommentCount(ctx context.Context, req *model.MomentCommentCount) error {
	return r.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"like_count":        gorm.Expr("like_count+VALUES(like_count)"),
			"like_cancel_count": gorm.Expr("like_cancel_count+VALUES(like_cancel_count)"),
			"comment_count":     gorm.Expr("comment_count+VALUES(comment_count)"),
		}),
	}).Create(req).Error
}

func NewCommunityRepository(r *Repository) CommunityRepository {
	return &communityRepository{
		Repository: r,
	}
}
