package ws

import (
	"errors"
	"fmt"
	"sync"
)

type ConnMgr struct {
	buckets      []*bucket
	perBucketCap int
}

func NewConnMgr(length, maxConns int) *ConnMgr {
	mgr := &ConnMgr{
		buckets:      make([]*bucket, length),
		perBucketCap: maxConns,
	}

	for i := 0; i < length; i++ {
		mgr.buckets[i] = NewBucket(i, maxConns)
	}

	return mgr
}

func (m *ConnMgr) AddConn(conn *WsConn) error {
	return m.GetBucket(conn.ConnId).Add(conn)
}

func (m *ConnMgr) GetConn(id int64) *WsConn {
	return m.GetBucket(id).Get(id)
}

func (m *ConnMgr) RemConn(id int64) error {
	return m.GetBucket(id).Rem(id)
}

func (m *ConnMgr) GetBucket(id int64) *bucket {
	index := id % int64(len(m.buckets))
	return m.buckets[index]
}

type bucket struct {
	mutx  sync.RWMutex
	index int //第几个桶
	len   int //最大连接数
	conns map[int64]*WsConn
}

func NewBucket(index, len int) *bucket {
	return &bucket{
		index: index,
		len:   len,
		conns: make(map[int64]*WsConn, len),
	}
}

func (b *bucket) Add(conn *WsConn) error {
	b.mutx.Lock()
	defer b.mutx.Unlock()
	if len(b.conns) < b.len {
		if old, ok := b.conns[conn.ConnId]; ok {
			old.Close()
		}
		b.conns[conn.ConnId] = conn
		return nil
	}
	return errors.New(fmt.Sprintf("bucket %v 连接数已满", b.index))
}

func (b *bucket) Get(id int64) *WsConn {
	b.mutx.RLock()
	defer b.mutx.RUnlock()
	return b.conns[id]
}

func (b *bucket) Rem(id int64) error {
	b.mutx.Lock()
	defer b.mutx.Unlock()
	if old, ok := b.conns[id]; ok {
		old.Close()
		delete(b.conns, id)
		return nil
	}
	return errors.New(fmt.Sprintf("conn %v is not found", id))
}
