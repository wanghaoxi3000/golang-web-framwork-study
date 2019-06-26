package ziface

// IRequest 把客户端请求的连接信息，和请求的数据包装到了一个 Request 中
type IRequest interface {
	GetConnection() IConnection // 得到当前链接
	GetData() []byte            // 得到请求的消息数据
	GetMsgID() uint32           // 得到请求的消息ID
}
