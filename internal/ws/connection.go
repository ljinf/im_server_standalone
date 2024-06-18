package ws

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/ljinf/im_server_standalone/pkg/log"
	"sync"
	"sync/atomic"
)

const (
	WriteChanMaxLen = 16
)

type WsConn struct {
	logger  *log.Logger
	ConnId  int64 //userId
	Conn    *websocket.Conn
	outChan chan []byte
	isClose int32 // 0否  1是
	once    sync.Once
}

func NewWsConn(logger *log.Logger, connId int64, conn *websocket.Conn) *WsConn {
	return &WsConn{
		logger:  logger,
		ConnId:  connId,
		Conn:    conn,
		outChan: make(chan []byte, WriteChanMaxLen),
	}
}

func (c *WsConn) Work(handler Dispatch) {
	c.logger.Debug(fmt.Sprintf("conn %v start read and write......", c.ConnId))
	go c.writeLoop()
	go c.readLoop(handler)
}

func (c *WsConn) readLoop(handler Dispatch) {
	for {
		messageType, payload, err := c.Conn.ReadMessage()
		if err != nil {
			c.logger.Error(err.Error())
			return
		}
		if messageType == websocket.PingMessage || messageType == websocket.PongMessage {
			continue
		}
		handler(c.ConnId, payload)
	}
}

func (c *WsConn) Write(payload []byte) error {
	if atomic.LoadInt32(&c.isClose) == 0 {
		c.outChan <- payload
		return nil
	}
	return errors.New("closed")
}

func (c *WsConn) writeLoop() {
	for v := range c.outChan {
		if err := c.Conn.WriteMessage(websocket.BinaryMessage, v); err != nil {
			c.logger.Error(err.Error())
			break
		}
	}
	c.Close()
}

func (c *WsConn) Close() {
	c.once.Do(func() {
		_ = c.Conn.Close()
		close(c.outChan)
		atomic.CompareAndSwapInt32(&c.isClose, 0, 1)
	})
}
