package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

//ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁
}

//NewConnManager 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//Add 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源 map， 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将 conn 加入到 ConnManager 中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager successfully: conn num = ", connMgr.Len())
}

//Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源map， 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("connID = ", conn.GetConnID(), " remove from ConnManager successfully: conn num = ", connMgr.Len())
}

//Get 根据connID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源map， 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

//Len 得到当前连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//ClearConn 清除并终止所有的连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源map， 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections succ!  conn num = ", connMgr.Len())
}
