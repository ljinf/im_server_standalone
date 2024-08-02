package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/ljinf/im_server_standalone/internal/cache"
	"github.com/ljinf/im_server_standalone/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
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
	DecrMsgSeq(ctx context.Context, convId int64)
	SelectConversationMsg(ctx context.Context, conversationId, seq int64, pageNum, pageSize int) ([]model.MsgResp, error)
	SelectLastConversationMsg(ctx context.Context, conversationId int64) (*model.MsgResp, error)

	// 用户消息链
	CreateUserMsgList(ctx context.Context, req *model.UserMsgList) error

	// 用户会话链
	CreateUserConversationList(ctx context.Context, req ...*model.UserConversationList) error
	UpdateUserConversationList(ctx context.Context, req *model.UserConversationList) error
	// 用户会话列表
	SelectUserConversationList(ctx context.Context, userId, pageNum, pageSize int64) ([]model.UserConversationList, error)
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

	if err = r.DB(ctx).Where("conversation_id in ?", conversationId).Find(&list).Error; err != nil {
		return nil, err
	}

	for _, v := range list {
		if err = cache.SetConversationCache(r.rdb, &v); err != nil {
			r.logger.Error(err.Error(), zap.Any("SetConversationCache", v))
		}
	}

	return list, nil
}

// 群聊的会话才会更新信息，例如公告等
func (r *chatRepository) UpdateConversation(ctx context.Context, req *model.ConversationList) error {
	return r.DB(ctx).Where("conversation_id=?", req.ConversationId).Updates(req).Error
}

func (r *chatRepository) CreateMsg(ctx context.Context, req *model.MsgList, seq int64) error {
	if err := r.DB(ctx).Create(req).Error; err != nil {
		return err
	}

	cacheInfo := &model.MsgResp{
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
	}

	//消息缓存
	if err := cache.SetMsgCache(r.rdb, cacheInfo); err != nil {
		r.logger.Error(err.Error(), zap.Any("SetMsgCache", cacheInfo))
	} else {
		//会话消息链缓存
		msgCount, err := cache.GetConversationMsgCount(r.rdb, cacheInfo.ConversationId)
		if err != nil {
			r.logger.Error(err.Error())
		}
		//队列溢出，删除部分以前的消息
		if int(msgCount) >= r.cacheMsgLength {
			if err = cache.RemConversationMsg(r.rdb, cacheInfo.ConversationId, int64(r.remCount)); err != nil {
				r.logger.Error(err.Error(), zap.Any("RemConversationMsg convId", cacheInfo.ConversationId))
			}
		} else {
			if err = cache.AddConversationMsgCache(r.rdb, *cacheInfo); err != nil {
				r.logger.Error(err.Error(), zap.Any("AddConversationMsgCache", cacheInfo))
			}
		}
	}
	return nil
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

// 如果数据库回滚，序号也要回滚
func (r *chatRepository) DecrMsgSeq(ctx context.Context, convId int64) {
	cache.DecrConversationMsg(r.rdb, convId)
}

func (r *chatRepository) SelectConversationMsg(ctx context.Context, conversationId, seq int64, pageNum, pageSize int) ([]model.MsgResp, error) {

	msgList, err := cache.GetConversationMsgList(r.rdb, conversationId, seq, int64(pageNum), int64(pageSize))
	if err != nil {
		r.logger.Error(err.Error(), zap.Any("conversationId", conversationId))
	}

	if len(msgList) >= pageSize {
		return msgList, nil
	}

	var list []model.MsgResp

	querySql := "SELECT cml.`seq`,ml.* FROM `conversation_msg_list` cml LEFT JOIN `msg_list` ml ON cml.`msg_id`=ml.`msg_id` " +
		"WHERE cml.`conversation_id`=? AND cml.`seq`>? ORDER BY cml.seq DESC LIMIT ? OFFSET ?"
	err = r.DB(ctx).Raw(querySql, conversationId, seq, pageSize, (pageNum-1)*pageSize).Scan(&list).Error
	return list, err
}

func (r *chatRepository) CreateUserMsgList(ctx context.Context, req *model.UserMsgList) error {
	return r.DB(ctx).Create(req).Error
}

// 创建会话信息
func (r *chatRepository) CreateUserConversationList(ctx context.Context, req ...*model.UserConversationList) error {
	if err := r.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"last_read_seq", "updated_at"}),
	}).Create(req).Error; err != nil {
		return err
	}

	for _, v := range req {
		if err := cache.SetUserConversationCache(r.rdb, *v); err != nil {
			r.logger.Error(err.Error(), zap.Any("SetUserConversationCache", req))
		}
	}
	return nil
}

// 更新会话信息
func (r *chatRepository) UpdateUserConversationList(ctx context.Context, req *model.UserConversationList) error {
	if err := r.DB(ctx).Where("user_id=? and conversation_id=?", req.UserId, req.ConversationId).Updates(req).Error; err != nil {
		return err
	}

	if userConversationCache, err := cache.GetUserConversationCache(r.rdb, req.UserId, req.ConversationId); err != nil {
		r.logger.Error(err.Error(), zap.Any("uid", req.UserId), zap.Any("convId", req.ConversationId))
	} else {
		userConversationCache.UserId = req.UserId
		userConversationCache.ConversationId = req.ConversationId
		if req.LastReadSeq != 0 {
			userConversationCache.LastReadSeq = req.LastReadSeq
		}
		if req.IsTop != 0 {
			userConversationCache.IsTop = req.IsTop
		}
		if req.NotifyType != 0 {
			userConversationCache.NotifyType = req.NotifyType
		}
		// 更新信息
		if err = cache.SetUserConversationCache(r.rdb, *userConversationCache); err != nil {
			r.logger.Error(err.Error())
		}
	}
	return nil
}

// 获取用户会话链信息
func (r *chatRepository) SelectUserConversationList(ctx context.Context, userId, pageNum, pageSize int64) ([]model.UserConversationList, error) {

	list, err := cache.GetUserConversationListCache(r.rdb, userId, pageNum, pageSize)
	if err != nil {
		r.logger.Error(err.Error(), zap.Any("uid", userId))
	}

	if len(list) > 0 {
		return list, nil
	}

	err = r.DB(ctx).Where("user_id=?", userId).Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&list).Error
	if len(list) > 0 {
		if err = cache.SetUserConversationCache(r.rdb, list...); err != nil {
			r.logger.Error(err.Error())
		}
	}

	return list, err
}

// 会话下的所有用户
func (r *chatRepository) SelectConversationUsers(ctx context.Context, conversationId int64) ([]model.UserInfo, error) {

	userInfoList, err := cache.GetConversationUserListCache(r.rdb, conversationId)
	if err != nil {
		r.logger.Error(err.Error(), zap.Any("conversationId", conversationId))
	}

	if len(userInfoList) > 0 {
		return userInfoList, nil
	}

	var info []model.AccountInfo
	querySql := "SELECT u.`user_id`,u.`nick_name`,u.`avatar`,u.`gender`,u.`status`,r.`email`,r.`phone` " +
		"FROM `user_info` u INNER JOIN `register` r ON u.`user_id`=r.`user_id` " +
		"WHERE u.`user_id` IN (SELECT uc.`user_id` FROM `user_conversation_list` uc WHERE uc.`conversation_id`=?)"
	if err = r.DB(ctx).Raw(querySql, conversationId).Scan(&info).Error; err != nil {
		return nil, err
	}

	if err = cache.SetAccountInfoCache(r.rdb, info...); err != nil {
		r.logger.Error(err.Error(), zap.Any("SetAccountInfoCache", info))
	} else {
		uids := make([]int64, 0, len(info))
		for _, v := range info {
			uids = append(uids, v.UserId)
		}
		if err = cache.AddConversationUserListCache(r.rdb, conversationId, uids...); err != nil {
			r.logger.Error(err.Error(), zap.Any("convId", conversationId),
				zap.Any("AddConversationUserListCache", uids))
		}
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

	if newestMsg, err := cache.GetConversationNewestMsg(r.rdb, conversationId); err != nil {
		r.logger.Error(err.Error(), zap.Any("GetConversationNewestMsg convId", conversationId))
	} else {
		return newestMsg, nil
	}

	//没有则从数据库加载一批最新的消息
	var list []model.MsgResp
	querySql := "SELECT cml.`seq`,ml.* FROM `conversation_msg_list` cml INNER JOIN `msg_list` ml ON cml.`msg_id`=ml.`msg_id` " +
		"WHERE cml.`conversation_id`=? ORDER BY cml.seq DESC LIMIT ? "
	if err := r.DB(ctx).Raw(querySql, conversationId, r.cacheMsgLength).Scan(&list).Error; err != nil {
		return nil, err
	}

	if len(list) > 0 {
		if err := cache.AddConversationMsgCache(r.rdb, list...); err != nil {
			r.logger.Error(err.Error())
		}
		return &list[0], nil
	}
	return nil, errors.New("ConversationNewestMsg Not Found")
}
