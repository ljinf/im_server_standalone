package repository

import (
	"context"
	"github.com/ljinf/im_server_standalone/internal/cache"
	"github.com/ljinf/im_server_standalone/internal/model"
	"gorm.io/gorm/clause"
)

type ChatRepository interface {
	// 会话
	CreateConversation(ctx context.Context, req *model.ConversationList) error
	SelectConversation(ctx context.Context, conversationId ...int64) ([]model.ConversationList, error)
	UpdateConversation(ctx context.Context, req *model.ConversationList) error

	// 消息
	CreateMsg(ctx context.Context, req *model.MsgList) error
	SelectMsgList(ctx context.Context, msgId ...int64) ([]model.MsgList, error)
	UpdateMsg(ctx context.Context, req *model.MsgList) error

	// 会话消息
	CreateConversationMsg(ctx context.Context, req *model.ConversationMsgList) error
	SelectConversationMsg(ctx context.Context, conversationId, seq int64, pageNum, pageSize int) ([]model.MsgResp, error)

	// 用户消息链
	CreateUserMsgList(ctx context.Context, req *model.UserMsgList) error

	// 用户会话链
	CreateUserConversationList(ctx context.Context, req ...*model.UserConversationList) error
	UpdateUserConversationList(ctx context.Context, req *model.UserConversationList) error
	// 用户会话列表
	SelectUserConversationList(ctx context.Context, userId int64) ([]model.ConversationResp, error)
	SelectConversationUsers(ctx context.Context, conversationId int64) ([]model.UserInfo, error) //会话下的用户列表
}

type chatRepository struct {
	*Repository
}

func NewChatRepository(r *Repository) ChatRepository {
	return &chatRepository{
		Repository: r,
	}
}

func (r *chatRepository) CreateConversation(ctx context.Context, req *model.ConversationList) error {
	return r.DB(ctx).Create(req).Error
}

func (r *chatRepository) SelectConversation(ctx context.Context, conversationId ...int64) ([]model.ConversationList, error) {
	var list []model.ConversationList
	if err := r.DB(ctx).Where("conversation_id in ?", conversationId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *chatRepository) UpdateConversation(ctx context.Context, req *model.ConversationList) error {
	return r.DB(ctx).Where("conversation_id=?", req.ConversationId).Updates(req).Error
}

func (r *chatRepository) CreateMsg(ctx context.Context, req *model.MsgList) error {
	return r.DB(ctx).Create(req).Error
}

func (r *chatRepository) SelectMsgList(ctx context.Context, msgId ...int64) ([]model.MsgList, error) {
	var list []model.MsgList
	if err := r.DB(ctx).Where("msg_id in ?", msgId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *chatRepository) UpdateMsg(ctx context.Context, req *model.MsgList) error {
	return r.DB(ctx).Where("msg_id=?", req.MsgId).Updates(req).Error
}

// 创建会话消息，并生成一个消息序列号
func (r *chatRepository) CreateConversationMsg(ctx context.Context, req *model.ConversationMsgList) error {
	msgSeq := cache.IncrConversationMsg(r.rdb, req.ConversationId)
	req.Seq = msgSeq
	return r.DB(ctx).Create(req).Error
}

func (r *chatRepository) SelectConversationMsg(ctx context.Context, conversationId, seq int64, pageNum, pageSize int) ([]model.MsgResp, error) {
	var list []model.MsgResp

	querySql := "SELECT cml.`seq`,ml.* FROM `conversation_msg_list` cml LEFT JOIN `msg_list` ml ON cml.`msg_id`=ml.`msg_id` " +
		"WHERE cml.`conversation_id`=? AND cml.`seq`>? ORDER BY cml.seq DESC LIMIT ? OFFSET ?"
	err := r.DB(ctx).Raw(querySql, conversationId, seq, pageSize, (pageNum-1)*pageSize).Scan(&list).Error
	return list, err
}

func (r *chatRepository) CreateUserMsgList(ctx context.Context, req *model.UserMsgList) error {
	return r.DB(ctx).Create(req).Error
}

func (r *chatRepository) CreateUserConversationList(ctx context.Context, req ...*model.UserConversationList) error {
	return r.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"last_read_seq", "updated_at"}),
	}).Create(req).Error
}

func (r *chatRepository) UpdateUserConversationList(ctx context.Context, req *model.UserConversationList) error {
	return r.DB(ctx).Where("user_id=? and conversation_id=?", req.UserId, req.ConversationId).Updates(req).Error
}

func (r *chatRepository) SelectUserConversationList(ctx context.Context, userId int64) ([]model.ConversationResp, error) {
	var list []model.ConversationResp
	querySql := "SELECT ucl.`user_id`,ucl.`conversation_id`,ucl.`is_top`,ucl.`last_read_seq`,ucl.`notify_type`,cl.`type`,cl.`avatar`,cml.`msg_id`,cml.`seq` " +
		"FROM `user_conversation_list` ucl INNER JOIN `conversation_list` cl ON ucl.`conversation_id`=cl.`conversation_id` " +
		"LEFT JOIN `conversation_msg_list` cml ON ucl.`conversation_id`=cml.`conversation_id` " +
		"WHERE ucl.`user_id`=?"
	err := r.DB(ctx).Raw(querySql, userId).Scan(&list).Error
	return list, err
}

// 会话下的所有用户
func (r *chatRepository) SelectConversationUsers(ctx context.Context, conversationId int64) ([]model.UserInfo, error) {
	querySql := "ELECT u.`user_id`,u.`avatar`,u.`nick_name` " +
		"FROM `user_info` u INNER JOIN `user_conversation_list` uc ON u.`user_id`=uc.`user_id` " +
		"WHERE uc.`conversation_id`=?"
	list := []model.UserInfo{}
	err := r.DB(ctx).Raw(querySql, conversationId).Find(&list).Error
	return list, err
}
