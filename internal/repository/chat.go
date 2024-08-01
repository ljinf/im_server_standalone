package repository

import (
	"context"
	"fmt"
	"github.com/ljinf/im_server_standalone/internal/cache"
	"github.com/ljinf/im_server_standalone/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"strconv"
)

type ChatRepository interface {
	// 会话
	CreateConversation(ctx context.Context, req *model.ConversationList) error
	SelectConversation(ctx context.Context, conversationId ...int64) ([]model.ConversationList, error)
	UpdateConversation(ctx context.Context, req *model.ConversationList) error

	// 消息
	CreateMsg(ctx context.Context, req *model.MsgList, seq int64) error
	SelectMsgList(ctx context.Context, msgId ...interface{}) ([]model.MsgResp, error)
	UpdateMsg(ctx context.Context, req *model.MsgList) error

	// 会话消息
	CreateConversationMsg(ctx context.Context, req *model.ConversationMsgList) error
	SelectConversationMsg(ctx context.Context, conversationId, seq int64, pageNum, pageSize int) ([]model.MsgResp, error)
	SelectLastConversationMsg(ctx context.Context, conversationId int64) (*model.MsgResp, error)

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
	if err := r.DB(ctx).Create(req).Error; err != nil {
		return err
	}
	if err := cache.SetConversationCache(r.rdb, req); err != nil {
		r.logger.Error(err.Error(), zap.Any("ConversationList", req))
	}
	return nil
}

func (r *chatRepository) SelectConversation(ctx context.Context, conversationId ...int64) ([]model.ConversationList, error) {
	var (
		list = make([]model.ConversationList, 0, len(conversationId))
	)

	conversationLists, err := cache.GetConversationCache(r.rdb, conversationId...)
	if err != nil {
		r.logger.Error(err.Error())
	}

	if len(conversationLists) > 0 {
		return conversationLists, nil
	}

	if err := r.DB(ctx).Where("conversation_id in ?", conversationId).Find(&list).Error; err != nil {
		return nil, err
	}

	for _, v := range list {
		if err := cache.SetConversationCache(r.rdb, &v); err != nil {
			r.logger.Error(err.Error(), zap.Any("SetConversationCache", v))
		}
	}

	return list, nil
}

// 群聊的会话才会更新信息
func (r *chatRepository) UpdateConversation(ctx context.Context, req *model.ConversationList) error {
	return r.DB(ctx).Where("conversation_id=?", req.ConversationId).Updates(req).Error
}

func (r *chatRepository) CreateMsg(ctx context.Context, req *model.MsgList, seq int64) error {
	if err := r.DB(ctx).Create(req).Error; err != nil {
		return err
	}

	return cache.SetMsgCache(r.rdb, &model.MsgResp{
		Id:             req.Id,
		UserId:         req.UserId,
		MsgId:          req.MsgId,
		ConversationId: req.ConversationId,
		Content:        req.Content,
		ContentType:    req.ContentType,
		Seq:            seq,
		Status:         req.Status,
		SendTime:       req.SendTime,
		CreatedAt:      req.CreatedAt,
	})
}

func (r *chatRepository) SelectMsgList(ctx context.Context, msgId ...interface{}) ([]model.MsgResp, error) {
	var (
		list []model.MsgResp
	)

	msgCache, err := cache.GetMsgCache(r.rdb, msgId...)
	if err != nil {
		r.logger.Error(err.Error(), zap.Any("msgids", msgId))
	}

	if len(msgCache) > 0 && len(msgCache) < len(msgId) {
		return msgCache, nil
	}

	//SELECT m.`msg_id`,m.`user_id`,m.`conversation_id`,m.`content`,m.`content_type`,m.`send_time`,m.`status`,cm.`seq`
	//FROM `msg_list` m INNER JOIN `conversation_msg_list` cm ON cm.`conversation_id`=m.`conversation_id WHERE m.`msg_id` IN ()
	if err = r.DB(ctx).Table("`msg_list` m").Select("m.`msg_id`,m.`user_id`,m.`conversation_id`,m.`content`,m.`content_type`,m.`send_time`,m.`status`,cm.`seq`").
		Joins("INNER JOIN `conversation_msg_list` cm on ON cm.`conversation_id`=m.`conversation_id").Where("m.`msg_id` IN ?", msgId).Find(&list).Error; err != nil {
		return nil, err
	}

	if err = cache.SetMsgCache(r.rdb, nil); err != nil {
		r.logger.Error(fmt.Sprintf("SetMsgCache %v", err))
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
	if err := r.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"last_read_seq", "updated_at"}),
	}).Create(req).Error; err != nil {
		return err
	}

	for _, v := range req {
		if err := cache.SetUserConversationCache(r.rdb, v); err != nil {
			r.logger.Error(err.Error(), zap.Any("SetUserConversationCache", req))
		}
	}
	return nil
}

func (r *chatRepository) UpdateUserConversationList(ctx context.Context, req *model.UserConversationList) error {
	if err := r.DB(ctx).Where("user_id=? and conversation_id=?", req.UserId, req.ConversationId).Updates(req).Error; err != nil {
		return err
	}
	if err := cache.DelUserConversationCache(r.rdb, req.UserId, strconv.Itoa(int(req.ConversationId))); err != nil {
		r.logger.Error(err.Error())
	}
	return nil
}

func (r *chatRepository) SelectUserConversationList(ctx context.Context, userId int64) ([]model.ConversationResp, error) {

	var list []model.ConversationResp
	querySql := "SELECT ucl.`user_id`,ucl.`conversation_id`,ucl.`is_top`,ucl.`last_read_seq`,ucl.`notify_type`,cl.`type`,cl.`avatar`" +
		"FROM `user_conversation_list` ucl INNER JOIN `conversation_list` cl ON ucl.`conversation_id`=cl.`conversation_id` " +
		"WHERE ucl.`user_id`=?"
	err := r.DB(ctx).Raw(querySql, userId).Scan(&list).Error
	return list, err
}

// 会话下的所有用户
func (r *chatRepository) SelectConversationUsers(ctx context.Context, conversationId int64) ([]model.UserInfo, error) {

	userIds, err := cache.GetConversationUserListCache(r.rdb, conversationId)
	if err != nil {
		r.logger.Error(err.Error(), zap.Any("conversationId", conversationId))
	}

	userInfoList, err := cache.GetUserInfoListCache(r.rdb, userIds...)
	if err != nil {
		r.logger.Error(err.Error(), zap.Any("uids", userIds))
	}

	if len(userInfoList) > 0 {
		return userInfoList, nil
	}

	var info []model.AccountInfo
	querySql := "SELECT u.`user_id`,u.`nick_name`,u.`avatar`,u.`gender`,u.`status`,r.`email`,r.`phone` " +
		"FROM `user_info` u INNER JOIN `register` r ON u.`user_id`=r.`user_id` " +
		"WHERE u.`user_id` IN (SELECT uc.`user_id` FROM `user_conversation_list` uc WHERE uc.`conversation_id`=?)"
	if err := r.DB(ctx).Raw(querySql, conversationId).Scan(&info).Error; err != nil {
		return nil, err
	}

	if err = cache.SetAccountInfoCache(r.rdb, info...); err != nil {
		r.logger.Error(err.Error(), zap.Any("info", info))
	}

	resp := make([]model.UserInfo, 0, len(info))
	for _, v := range info {
		resp = append(resp, model.UserInfo{
			UserId:   v.UserId,
			NickName: v.NickName,
			Avatar:   v.Avatar,
			Gender:   v.Gender,
			Status:   v.Status,
		})
	}

	return resp, nil
}

// 会话最新一条消息
func (r *chatRepository) SelectLastConversationMsg(ctx context.Context, conversationId int64) (*model.MsgResp, error) {
	var msg model.MsgResp

	querySql := "SELECT cm.`conversation_id`,cm.`msg_id`,cm.`seq`,m.`user_id`,m.`content`,m.`content_type`,m.`send_time`,m.`status` FROM `conversation_msg_list` cm " +
		"left join `msg_list` m on m.`msg_id`=cm.`msg_id` " +
		"WHERE cm.`conversation_id` =? ORDER BY cm.`seq` DESC LIMIT 1"

	err := r.DB(ctx).Raw(querySql, conversationId).Scan(&msg).Error
	return &msg, err
}
