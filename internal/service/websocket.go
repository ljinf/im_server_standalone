package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/ws"
	"github.com/ljinf/im_server_standalone/pkg/contants"
	"go.uber.org/zap"
)

type WebsocketService interface {
	InitConn(userId int64, conn *websocket.Conn)
	PushMsg(payload []byte, userIds ...int64)
	ProcessMsg(sender int64, payload []byte)
}

type websocketService struct {
	*Service
	ws.SocketWsServer
}

func NewWebsocketService(s *Service, wss ws.SocketWsServer) WebsocketService {
	return &websocketService{
		Service:        s,
		SocketWsServer: wss,
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

func (w *websocketService) PushMsg(payload []byte, userIds ...int64) {
	if err := w.Push(payload, userIds...); err != nil {
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

		break
	}
}
