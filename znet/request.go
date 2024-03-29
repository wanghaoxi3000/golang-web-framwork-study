package znet

import (
	"zinx/ziface"
)

type Request struct {
	conn ziface.IConnection // 已经和客户端建立好的连接
	msg  ziface.IMessage    // 客户端请求的数据
}

// GetConnection 得到当前链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 得到请求的消息ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
