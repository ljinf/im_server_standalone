package ws

import (
	"github.com/ljinf/im_server_standalone/pkg/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

type Dispatch func(sender int64, payload []byte)

var (
	server SocketWsServer
	mutex  sync.Mutex
)

type SocketWsServer interface {
	AddConn(c *WsConn) error
	GetConnManager() *ConnMgr
	Push(msg []byte, ids ...int64) error
}

type wsServer struct {
	logger  *log.Logger
	connMgr *ConnMgr
}

func NewWsServer(conf *viper.Viper, logger *log.Logger) SocketWsServer {
	if server == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if server == nil {
			server = &wsServer{
				logger:  logger,
				connMgr: NewConnMgr(conf.GetInt("ws_server.max_buckets"), conf.GetInt("ws_server.per_bucket_cap")),
			}
		}
		return server
	}
	return server
}

func (s *wsServer) GetConnManager() *ConnMgr {
	return s.connMgr
}

func (s *wsServer) AddConn(c *WsConn) error {
	return s.connMgr.AddConn(c)
}

func (s *wsServer) Push(msg []byte, ids ...int64) error {
	if len(ids) > 0 {
		for _, v := range ids {
			if wsConn := s.connMgr.GetConn(v); wsConn != nil {
				if err := wsConn.Write(msg); err != nil {
					s.logger.Error(err.Error(), zap.Any("userId", v), zap.Any("msg", string(msg)))
				}
			}
		}
	}
	return nil
}
