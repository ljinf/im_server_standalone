package service

import (
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/repository"
	"go.uber.org/zap"
	"time"
)

type ChatService interface {
	CreateMsg(ctx context.Context, req *v1.SendMsgReq) (*v1.SendMsgResp, error)
	//历史消息
	GetMsgList(ctx context.Context, userId, conversationId string, seq int64, pageNum, pageSize int) ([]v1.SendMsgResp, error)

	// 会话
	GetUserConversationList(ctx context.Context, userId string, pageNum, pageSize int64) ([]v1.ConversationResp, error)
	GetConversationUsers(ctx context.Context, conversationId string) ([]v1.GetProfileResponseData, error) //会话下的用户
	//创建会话
	CreateConversationList(ctx context.Context, list ...*model.ConversationList) error

	//该会话最新一条消息
	GetLastConversationMsg(ctx context.Context, conversationId string) v1.SendMsgResp

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

	var (
		now = time.Now().Unix()

		// 消息序列号
		mSeq int64
		//消息在会话中的序列号，用于保证消息的顺序
		cSeq  int64
		msgId string
	)

	id, err := s.sid.GenUint64()
	if err != nil {
		return nil, err
	}
	msgId = fmt.Sprintf("%v", id)

	msg := &model.MsgList{
		UserId:         req.UserId,
		MsgId:          msgId,
		ConversationId: req.ConversationId,
		Content:        req.Content,
		ContentType:    req.ContentType,
		Seq:            cSeq,
		Status:         0,
		SentAt:         now,
	}

	//会话不存在
	if len(msg.ConversationId) == 0 {
		if err = s.tm.Transaction(ctx, func(ctx context.Context) error {
			// 用户会话链
			userConversationList := make([]*model.UserConversationList, 0, 2)

			convId := []string{req.UserId, req.TargetId}
			slice.Sort(convId)
			// 单聊，如果是群聊，要显式创建，所以会话id不为空
			conversationInfo := &model.ConversationList{
				ConversationId: fmt.Sprintf("%v-%v", convId[0], convId[1]),
				Type:           0,
				Member:         2,
				RecentMsgTime:  now,
				CreatedAt:      now,
			}
			//创建会话
			if err = s.repo.CreateConversation(ctx, conversationInfo); err != nil {
				return err
			}
			msg.ConversationId = conversationInfo.ConversationId
			// 接收者的会话
			userConversationList = append(userConversationList, &model.UserConversationList{
				UserId:         req.TargetId,
				ConversationId: msg.ConversationId,
				CreatedAt:      now,
			})

			// 发送者的会话链
			userConversationList = append(userConversationList, &model.UserConversationList{
				UserId:         req.UserId,
				ConversationId: msg.ConversationId,
				LastReadSeq:    cSeq,
				CreatedAt:      now,
			})
			if err = s.repo.CreateUserConversationList(ctx, userConversationList...); err != nil {
				return err
			}
			return nil
		}); err != nil {
			s.logger.Error(err.Error(), zap.Any("create conversation", ""))
			return nil, err
		}
	}

	if err = s.tm.Transaction(ctx, func(ctx context.Context) error {

		//用户消息链
		userMsgList := []*model.UserMsgList{
			&model.UserMsgList{
				UserId:         req.UserId,
				MsgId:          msgId,
				ConversationId: msg.ConversationId,
				Seq:            mSeq,
			},
			&model.UserMsgList{
				UserId:         req.TargetId,
				MsgId:          msgId,
				ConversationId: msg.ConversationId,
				Seq:            mSeq,
			}}
		if err = s.repo.CreateUserMsgList(ctx, userMsgList...); err != nil {
			return err
		}

		//消息体
		return s.repo.CreateMsg(ctx, msg, mSeq)
	}); err != nil {
		s.logger.Error(err.Error(), zap.Any("req", req))
		return nil, err
	}

	resp := &v1.SendMsgResp{
		UserId:         msg.UserId,
		MsgId:          msg.MsgId,
		ConversationId: msg.ConversationId,
		Content:        msg.Content,
		ContentType:    msg.ContentType,
		Status:         msg.Status,
		Seq:            msg.Seq,
		SendTime:       msg.SentAt,
	}
	return resp, nil
}

func (s *chatService) GetMsgList(ctx context.Context, userId, conversationId string, seq int64, pageNum, pageSize int) ([]v1.SendMsgResp, error) {
	/*msgLists, err := s.repo.SelectConversationMsg(ctx, conversationId, seq, pageNum, pageSize)
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
	return resp, nil*/
	return nil, nil
}

func (s *chatService) GetUserConversationList(ctx context.Context, userId string, pageNum, pageSize int64) ([]v1.ConversationResp, error) {
	/*userConversationList, err := s.repo.SelectUserConversationList(ctx, userId, pageNum, pageSize)
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
	return resp, nil*/
	return nil, nil
}

func (s *chatService) CreateConversationList(ctx context.Context, list ...*model.ConversationList) error {

	//TODO implement me
	panic("implement me")
}

func (s *chatService) ReportReadMsgSeq(ctx context.Context, req *v1.ReportReadReq) error {
	//now := time.Now().Unix()
	uc := model.UserConversationList{
		UserId:         req.UserId,
		ConversationId: req.ConversationId,
		LastReadSeq:    req.Seq,
	}

	return s.repo.UpdateUserConversationList(ctx, &uc)
}

func (s *chatService) GetConversationUsers(ctx context.Context, conversationId string) ([]v1.GetProfileResponseData, error) {
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

func (s *chatService) GetLastConversationMsg(ctx context.Context, conversationId string) v1.SendMsgResp {
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
		//SendTime:       lastMsg.SendTime,
	}
}
