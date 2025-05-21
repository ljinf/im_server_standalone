package ws

import (
	"errors"
	"fmt"
	"strconv"
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

func (m *ConnMgr) GetConn(id string) *WsConn {
	return m.GetBucket(id).Get(id)
}

func (m *ConnMgr) RemConn(id string) error {
	return m.GetBucket(id).Rem(id)
}

func (m *ConnMgr) GetBucket(id string) *bucket {
	atoi, _ := strconv.Atoi(id)
	index := atoi % len(m.buckets)
	return m.buckets[index]
}

type bucket struct {
	mutx  sync.RWMutex
	index int //第几个桶
	len   int //最大连接数
	conns map[string]*WsConn
}

func NewBucket(index, len int) *bucket {
	return &bucket{
		index: index,
		len:   len,
		conns: make(map[string]*WsConn, len),
	}
}

func (b *bucket) Add(conn *WsConn) error {
	if len(b.conns) < b.len {
		b.mutx.Lock()
		old, ok := b.conns[conn.ConnId]
		b.conns[conn.ConnId] = conn
		b.mutx.Unlock()
		if ok {
			old.Close()
		}
		return nil
	}
	return errors.New(fmt.Sprintf("bucket %v 连接数已满", b.index))
}

func (b *bucket) Get(id string) *WsConn {
	b.mutx.RLock()
	defer b.mutx.RUnlock()
	return b.conns[id]
}

func (b *bucket) Rem(id string) error {
	b.mutx.Lock()
	defer b.mutx.Unlock()
	if _, ok := b.conns[id]; ok {
		delete(b.conns, id)
		return nil
	}
	return errors.New(fmt.Sprintf("conn %v is not found", id))
}
