package service

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/ws"
	"github.com/ljinf/im_server_standalone/pkg/contants"
	"github.com/panjf2000/ants"
	"go.uber.org/zap"
)

type WebsocketService interface {
	InitConn(userId int64, conn *websocket.Conn)
	PushMsg(payload []byte, userIds ...int64)
	SyncPushMsg(msgInfo interface{}, userIds ...int64)
	ProcessMsg(sender int64, payload []byte)
}

type websocketService struct {
	*Service
	ws.SocketWsServer
	chatSrv ChatService
	task    *ants.Pool
}

func NewWebsocketService(s *Service, wss ws.SocketWsServer, chatSrv ChatService, pool *ants.Pool) WebsocketService {
	return &websocketService{
		Service:        s,
		SocketWsServer: wss,
		chatSrv:        chatSrv,
		task:           pool,
	}
}

func (w *websocketService) InitConn(userId int64, conn *websocket.Conn) {
	wsConn := ws.NewWsConn(w.logger, userId, conn)
	if err := w.AddConn(wsConn); err != nil {
		w.logger.Error(err.Error(), zap.Any("userId", userId))
		wsConn.Close()
		return
	}

	wsConn.Work(w.ProcessMsg)
}

// 推送
func (w *websocketService) PushMsg(payload []byte, userIds ...int64) {
	if err := w.Push(payload, userIds...); err != nil {
		w.logger.Error(err.Error())
	}
}

// 移步推送
func (w *websocketService) SyncPushMsg(msgInfo interface{}, userIds ...int64) {
	if err := w.task.Submit(func() {
		payload, err := json.Marshal(msgInfo)
		if err != nil {
			w.logger.Error(err.Error())
			return
		}
		w.PushMsg(payload, userIds...)
	}); err != nil {
		w.logger.Error(err.Error())
	}
}

// 消息处理
func (w *websocketService) ProcessMsg(sender int64, payload []byte) {
	var info model.WsMessage
	if err := json.Unmarshal(payload, &info); err != nil {
		w.logger.Error(err.Error(), zap.Any("消息内容解析错误 payload", string(payload)))
		return
	}

	switch info.MsgType {
	case contants.MsgTypeCommand:

		break
	case contants.MsgTypeNotify:

		break
	case contants.MsgTypeChat:
		w.msgChat(info.Payload)
		break
	}
}

func (w *websocketService) msgChat(payload []byte) {
	msgReq, err := parsePayload(payload)
	if err != nil {
		w.logger.Error(err.Error(), zap.Any("msgChat", "parsePayload err"))
		return
	}

	msgResp, err := w.chatSrv.CreateMsg(context.Background(), msgReq)
	if err != nil {
		w.logger.Error(err.Error(), zap.Any("msgChat", "CreateMsg err"))
	}

	resp, err := json.Marshal(msgResp)
	if err != nil {
		w.logger.Error(err.Error(), zap.Any("msgChat", "json.Marshal err"))
		return
	}

	w.PushMsg(resp, msgResp.UserId, msgReq.TargetId)
}

func parsePayload(payload []byte) (*v1.SendMsgReq, error) {
	msg := &v1.SendMsgReq{}
	if err := json.Unmarshal(payload, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
