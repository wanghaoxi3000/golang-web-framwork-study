package ziface

import "net"

// IConnection 链接模块接口
type IConnection interface {
	Start()                         // 启动连接
	Stop()                          // 停止连接
	GetTCPConnection() *net.TCPConn // 获取当前连接绑定的 socket conn
	GetConnID() uint32              // 获取当前链接模块的链接ID
	RemoteAddr() net.Addr           // 获取远程客户端 TCP 状态 IP port
	SendMsg(uint32, []byte) error   // 发送数据，将数据发送到远程的客户端

	SetProperty(key string, value interface{})   // 设置链接属性
	GetProperty(key string) (interface{}, error) // 获取链接属性
	RemoveProperty(key string)                   // 移除链接属性
}

// HandleFunc 处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
