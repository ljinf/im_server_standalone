package service

import (
	"context"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/repository"
	"github.com/ljinf/im_server_standalone/pkg/contants"
	"go.uber.org/zap"
	"time"
)

type ChatService interface {
	CreateMsg(ctx context.Context, req *v1.SendMsgReq) (*v1.SendMsgResp, error)
	//历史消息
	GetMsgList(ctx context.Context, userId, conversationId, seq int64, pageNum, pageSize int) ([]v1.SendMsgResp, error)

	// 会话
	GetUserConversationList(ctx context.Context, userId, pageNum, pageSize int64) ([]v1.ConversationResp, error)
	GetConversationUsers(ctx context.Context, conversationId int64) ([]v1.GetProfileResponseData, error) //会话下的用户
	//创建会话
	CreateConversationList(ctx context.Context, list ...*model.ConversationList) error

	//该会话最新一条消息
	GetLastConversationMsg(ctx context.Context, conversationId int64) v1.SendMsgResp

	//已读上报
	ReportReadMsgSeq(ctx context.Context, req *v1.ReportReadReq) error
}

type chatService struct {
	*Service
	repo repository.ChatRepository
}

func NewChatService(s *Service, repo repository.ChatRepository) ChatService {
	return &chatService{
		Service: s,
		repo:    repo,
	}
}

// 返回消息ID
func (s *chatService) CreateMsg(ctx context.Context, req *v1.SendMsgReq) (*v1.SendMsgResp, error) {

	msgId, err := s.sid.GenUint64()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	msg := &model.MsgList{
		UserId:         req.UserId,
		MsgId:          int64(msgId),
		ConversationId: req.ConversationId,
		Content:        req.Content,
		ContentType:    req.ContentType,
		Status:         0,
		SendTime:       now,
		CreatedAt:      now,
	}

	// 消息序列号
	var mSeq int64

	if err = s.tm.Transaction(ctx, func(ctx context.Context) error {

		// 用户会话链
		userConversationList := make([]*model.UserConversationList, 0, 2)
		//会话不存在
		if msg.ConversationId == 0 {
			cId, err := s.sid.GenUint64()
			if err != nil {
				return err
			}
			// 单聊，如果是群聊，要先创建群，所以会话id不为0
			conversationInfo := &model.ConversationList{
				ConversationId: int64(cId),
				Type:           0,
				Member:         2,
				RecentMsgTime:  now,
				CreatedAt:      now,
			}
			if err = s.repo.CreateConversation(ctx, conversationInfo); err != nil {
				return err
			}
			msg.ConversationId = conversationInfo.ConversationId
			// 接收者的会话
			userConversationList = append(userConversationList, &model.UserConversationList{
				UserId:         req.TargetId,
				ConversationId: msg.ConversationId,
				CreatedAt:      now,
				UpdatedAt:      now,
			})
		}

		//会话消息
		cMsg := &model.ConversationMsgList{
			ConversationId: msg.ConversationId,
			MsgId:          msg.MsgId,
			CreatedAt:      now,
		}
		if err = s.repo.CreateConversationMsg(ctx, cMsg); err != nil {
			return err
		}

		// 发送者的会话链
		userConversationList = append(userConversationList, &model.UserConversationList{
			UserId:         req.UserId,
			ConversationId: msg.ConversationId,
			LastReadSeq:    cMsg.Seq,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
		if err = s.repo.CreateUserConversationList(ctx, userConversationList...); err != nil {
			return err
		}
		mSeq = cMsg.Seq
		//消息体
		return s.repo.CreateMsg(ctx, msg, mSeq)
	}); err != nil {
		if mSeq > 0 {
			//序号回滚
			s.repo.DecrMsgSeq(ctx, msg.ConversationId)
		}
		s.logger.Error(err.Error(), zap.Any("req", req))
		return nil, err
	}

	resp := &v1.SendMsgResp{
		UserId:         msg.UserId,
		MsgId:          int64(msgId),
		ConversationId: msg.ConversationId,
		Content:        msg.Content,
		ContentType:    msg.ContentType,
		Status:         msg.Status,
		Seq:            mSeq,
		SendTime:       msg.SendTime,
		CreatedAt:      now,
	}
	return resp, nil
}

func (s *chatService) GetMsgList(ctx context.Context, userId, conversationId, seq int64, pageNum, pageSize int) ([]v1.SendMsgResp, error) {
	msgLists, err := s.repo.SelectConversationMsg(ctx, conversationId, seq, pageNum, pageSize)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("conversationId", conversationId))
		return nil, v1.ErrInternalServerError
	}

	resp := make([]v1.SendMsgResp, 0, len(msgLists))
	for _, v := range msgLists {
		resp = append(resp, v1.SendMsgResp{
			UserId:         v.UserId,
			MsgId:          v.MsgId,
			ConversationId: v.ConversationId,
			Content:        v.Content,
			ContentType:    v.ContentType,
			Status:         v.Status,
			Seq:            v.Seq,
			SendTime:       v.SendTime,
		})
	}
	return resp, nil
}

func (s *chatService) GetUserConversationList(ctx context.Context, userId, pageNum, pageSize int64) ([]v1.ConversationResp, error) {
	userConversationList, err := s.repo.SelectUserConversationList(ctx, userId, pageNum, pageSize)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("userId", userId))
		return nil, v1.ErrInternalServerError
	}

	convIds := make([]int64, 0, len(userConversationList))
	for _, v := range userConversationList {
		convIds = append(convIds, v.ConversationId)
	}
	conversationLists, err := s.repo.SelectConversation(ctx, convIds...)
	if err != nil {
		return nil, v1.ErrInternalServerError
	}

	resp := make([]v1.ConversationResp, 0, len(userConversationList))
	for index, v := range userConversationList {
		conv := v1.ConversationResp{
			ConversationId: v.ConversationId,
			Type:           conversationLists[index].Type,
			Avatar:         conversationLists[index].Avatar,
			LastReadSeq:    v.LastReadSeq,
			NotifyType:     v.NotifyType,
			IsTop:          v.IsTop,
			RecentMsg:      s.GetLastConversationMsg(ctx, v.ConversationId),
		}
		//单聊会话获取用户列表
		if conv.Type == contants.ConversationTypeC2C {
			conversationUsers, _ := s.GetConversationUsers(ctx, v.ConversationId)
			conv.UserList = conversationUsers //会话用户列表
		}
		resp = append(resp, conv)
	}
	return resp, nil
}

func (s *chatService) CreateConversationList(ctx context.Context, list ...*model.ConversationList) error {

	//TODO implement me
	panic("implement me")
}

func (s *chatService) ReportReadMsgSeq(ctx context.Context, req *v1.ReportReadReq) error {
	now := time.Now().Unix()
	uc := model.UserConversationList{
		UserId:         req.UserId,
		ConversationId: req.ConversationId,
		LastReadSeq:    req.Seq,
		UpdatedAt:      now,
	}

	return s.repo.UpdateUserConversationList(ctx, &uc)
}

func (s *chatService) GetConversationUsers(ctx context.Context, conversationId int64) ([]v1.GetProfileResponseData, error) {
	users, err := s.repo.SelectConversationUsers(ctx, conversationId)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("convId", conversationId))
	}

	resp := make([]v1.GetProfileResponseData, 0, len(users))

	for _, v := range users {
		user := v1.GetProfileResponseData{
			UserId:   v.UserId,
			Avatar:   v.Avatar,
			NickName: v.NickName,
		}
		resp = append(resp, user)
	}

	return resp, nil
}

func (s *chatService) GetLastConversationMsg(ctx context.Context, conversationId int64) v1.SendMsgResp {
	lastMsg, err := s.repo.SelectLastConversationMsg(ctx, conversationId)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("convId", conversationId))
	}
	return v1.SendMsgResp{
		ConversationId: lastMsg.ConversationId,
		MsgId:          lastMsg.MsgId,
		UserId:         lastMsg.UserId,
		Seq:            lastMsg.Seq,
		Content:        lastMsg.Content,
		ContentType:    lastMsg.ContentType,
		Status:         lastMsg.Status,
		SendTime:       lastMsg.SendTime,
	}
}
