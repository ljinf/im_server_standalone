package service

import (
	"context"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/repository"
	"go.uber.org/zap"
	"time"
)

type ChatService interface {
	CreateMsg(ctx context.Context, req *v1.MsgReq) (*v1.MsgResp, error)
	GetMsgList(ctx context.Context, userId, conversationId, seq int64, pageNum, pageSize int) ([]v1.MsgResp, error)

	GetUserConversationList(ctx context.Context, userId int64) ([]v1.ConversationResp, error)
	CreateConversationList(ctx context.Context, list ...*model.ConversationList) error
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
func (s *chatService) CreateMsg(ctx context.Context, req *v1.MsgReq) (*v1.MsgResp, error) {

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
		SendTime:       req.SendTime,
		CreatedAt:      now,
	}

	// 消息序列号
	var mSeq int64

	if err = s.tm.Transaction(ctx, func(ctx context.Context) error {

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

		// 用户会话链
		ucl := []*model.UserConversationList{
			{
				UserId:         req.UserId,
				ConversationId: msg.ConversationId,
				LastReadSeq:    cMsg.Seq,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			{
				UserId:         req.TargetId,
				ConversationId: msg.ConversationId,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		}
		if err = s.repo.CreateUserConversationList(ctx, ucl...); err != nil {
			return err
		}
		mSeq = cMsg.Seq
		//消息体
		return s.repo.CreateMsg(ctx, msg)
	}); err != nil {
		s.logger.Error(err.Error(), zap.Any("req", req))
		return nil, err
	}

	resp := &v1.MsgResp{
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

func (s *chatService) GetMsgList(ctx context.Context, userId, conversationId, seq int64, pageNum, pageSize int) ([]v1.MsgResp, error) {
	msgLists, err := s.repo.SelectConversationMsg(ctx, conversationId, seq, pageNum, pageSize)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("conversationId", conversationId))
		return nil, v1.ErrInternalServerError
	}

	resp := make([]v1.MsgResp, 0, len(msgLists))
	for _, v := range msgLists {
		resp = append(resp, v1.MsgResp{
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

func (s *chatService) GetUserConversationList(ctx context.Context, userId int64) ([]v1.ConversationResp, error) {
	userConversationList, err := s.repo.SelectUserConversationList(ctx, userId)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("userId", userId))
		return nil, v1.ErrInternalServerError
	}

	resp := make([]v1.ConversationResp, 0, len(userConversationList))
	for _, v := range userConversationList {
		resp = append(resp, v1.ConversationResp{
			ConversationId: v.ConversationId,
			Type:           v.Type,
			Avatar:         v.Avatar,
			RecentMsg: v1.MsgResp{
				MsgId: v.MsgId,
				Seq:   v.Seq,
			}, //TODO
			LastReadSeq: v.LastReadSeq,
			NotifyType:  v.NotifyType,
			IsTop:       v.IsTop,
		})
	}
	return resp, nil
}

func (s *chatService) CreateConversationList(ctx context.Context, list ...*model.ConversationList) error {

	//TODO implement me
	panic("implement me")
}
