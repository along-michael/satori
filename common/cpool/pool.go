package cpool

import (
	"fmt"
	"sync"
	"time"
)

//TODO: 保存所有的连接, 而不是只保存连接计数

var ErrMaxConn = fmt.Errorf("maximum connections reached")

//
type ConnPool struct {
	sync.RWMutex

	Name        string
	Address     string
	MaxConns    int
	MaxIdle     int
	ConnTimeout int
	CallTimeout int
	Cnt         int64
	New         func(name string, pool *ConnPool) (PoolClient, error)

	active int
	free   []PoolClient
	all    map[string]PoolClient
}

type ConnPoolStats struct {
	Name   string
	Count  int64
	Active int
	All    int
	Free   int
}

func (this *ConnPoolStats) String() string {
	return fmt.Sprintf("%s[Count: %d, Active: %d, All: %d, Free: %d]",
		this.Name,
		this.Count,
		this.Active,
		this.All,
		this.Free,
	)
}

func NewConnPool(name string, address string, maxConns int, maxIdle int, connTimeout int, callTimeout int, new func(string, *ConnPool) (PoolClient, error)) *ConnPool {
	return &ConnPool{
		Name:        name,
		Address:     address,
		MaxConns:    maxConns,
		MaxIdle:     maxIdle,
		CallTimeout: callTimeout,
		ConnTimeout: connTimeout,
		Cnt:         0,
		New:         new,
		all:         make(map[string]PoolClient),
	}
}

func (this *ConnPool) Stats() *ConnPoolStats {
	this.RLock()
	defer this.RUnlock()

	return &ConnPoolStats{
		Name:   this.Name,
		Count:  this.Cnt,
		Active: this.active,
		All:    len(this.all),
		Free:   len(this.free),
	}

}

func (this *ConnPool) Fetch() (PoolClient, error) {
	this.Lock()
	defer this.Unlock()

	// get from free
	conn := this.fetchFree()
	if conn != nil {
		return conn, nil
	}

	if this.active >= this.MaxConns {
		return nil, ErrMaxConn
	}

	// create new conn
	conn, err := this.newConn()
	if err != nil {
		return nil, err
	}

	this.active += 1
	return conn, nil
}

func (this *ConnPool) Call(arg interface{}) (interface{}, error) {
	conn, err := this.Fetch()
	if err != nil {
		return nil, fmt.Errorf("%s get connection fail: conn %v, err %v. stats: %s", this.Name, conn, err, this.Stats())
	}

	callTimeout := time.Duration(this.CallTimeout) * time.Millisecond

	done := make(chan error)
	var resp interface{}
	go func() {
		resp, err = conn.Call(arg)
		done <- err
	}()

	select {
	case <-time.After(callTimeout):
		this.ForceClose(conn)
		return nil, fmt.Errorf("%s, call timeout", conn.Name())
	case err = <-done:
		if err != nil {
			this.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. stats: %s", this.Name, err, this.Stats())
		} else {
			this.Release(conn)
		}
		return resp, err
	}
}

func (this *ConnPool) Release(conn PoolClient) {
	this.Lock()
	defer this.Unlock()

	if len(this.free) >= this.MaxIdle {
		this.deleteConn(conn)
		this.active -= 1
	} else {
		this.addFree(conn)
	}
}

func (this *ConnPool) ForceClose(conn PoolClient) {
	this.Lock()
	defer this.Unlock()

	this.deleteConn(conn)
	this.active -= 1
}

func (this *ConnPool) Destroy() {
	this.Lock()
	defer this.Unlock()

	for _, conn := range this.free {
		if conn != nil && !conn.Closed() {
			conn.Close()
		}
	}

	for _, conn := range this.all {
		if conn != nil && !conn.Closed() {
			conn.Close()
		}
	}

	this.active = 0
	this.free = []PoolClient{}
	this.all = map[string]PoolClient{}
}

// internal, concurrently unsafe
func (this *ConnPool) newConn() (PoolClient, error) {
	name := fmt.Sprintf("%s_%d_%d", this.Name, this.Cnt, time.Now().Unix())
	conn, err := this.New(name, this)
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		return nil, err
	}

	this.Cnt++
	this.all[conn.Name()] = conn
	return conn, nil
}

func (this *ConnPool) deleteConn(conn PoolClient) {
	if conn != nil {
		conn.Close()
	}
	delete(this.all, conn.Name())
}

func (this *ConnPool) addFree(conn PoolClient) {
	this.free = append(this.free, conn)
}

func (this *ConnPool) fetchFree() PoolClient {
	if len(this.free) == 0 {
		return nil
	}

	conn := this.free[0]
	this.free = this.free[1:]
	return conn
}
