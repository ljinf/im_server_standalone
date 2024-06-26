package repository

import (
	"context"
	"github.com/ljinf/im_server_standalone/internal/cache"
	"github.com/ljinf/im_server_standalone/internal/model"
)

type IMMessageRepository interface {
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

	// 用户消息链
	CreateUserMsgList(ctx context.Context, req *model.UserMsgList) error

	// 用户会话链
	CreateUserConversationList(ctx context.Context, req ...*model.UserConversationList) error
	UpdateUserConversationList(ctx context.Context, req *model.UserConversationList) error
}

type imMessageRepository struct {
	*Repository
}

func NewIMMessageRepository(r *Repository) IMMessageRepository {
	return &imMessageRepository{
		Repository: r,
	}
}

func (r *imMessageRepository) CreateConversation(ctx context.Context, req *model.ConversationList) error {
	return r.DB(ctx).Create(req).Error
}

func (r *imMessageRepository) SelectConversation(ctx context.Context, conversationId ...int64) ([]model.ConversationList, error) {
	var list []model.ConversationList
	if err := r.DB(ctx).Where("conversation_id in ?", conversationId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *imMessageRepository) UpdateConversation(ctx context.Context, req *model.ConversationList) error {
	return r.DB(ctx).Where("conversation_id=?", req.ConversationId).Updates(req).Error
}

func (r *imMessageRepository) CreateMsg(ctx context.Context, req *model.MsgList) error {
	return r.DB(ctx).Create(req).Error
}

func (r *imMessageRepository) SelectMsgList(ctx context.Context, msgId ...int64) ([]model.MsgList, error) {
	var list []model.MsgList
	if err := r.DB(ctx).Where("msg_id in ?", msgId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *imMessageRepository) UpdateMsg(ctx context.Context, req *model.MsgList) error {
	return r.DB(ctx).Where("msg_id=?", req.MsgId).Updates(req).Error
}

// 创建会话消息，并生成一个消息序列号
func (r *imMessageRepository) CreateConversationMsg(ctx context.Context, req *model.ConversationMsgList) error {
	msgSeq := cache.IncrConversationMsg(r.rdb, req.ConversationId)
	req.Seq = msgSeq
	return r.DB(ctx).Create(req).Error
}

func (r *imMessageRepository) CreateUserMsgList(ctx context.Context, req *model.UserMsgList) error {
	return r.DB(ctx).Create(req).Error
}

func (r *imMessageRepository) CreateUserConversationList(ctx context.Context, req ...*model.UserConversationList) error {
	return r.DB(ctx).Create(req).Error
}

func (r *imMessageRepository) UpdateUserConversationList(ctx context.Context, req *model.UserConversationList) error {
	return r.DB(ctx).Where("user_id=? and conversation_id=?", req.UserId, req.ConversationId).Updates(req).Error
}
