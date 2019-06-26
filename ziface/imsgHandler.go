package ziface

//IMsgHandle 消息管理抽象层
type IMsgHandle interface {
	DoMsgHandler(request IRequest)          // 调度/执行对应的 Router 消息处理方法
	AddRouter(msgID uint32, router IRouter) // 为消息添加具体的处理逻辑
}
