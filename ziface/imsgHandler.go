package ziface

//IMsgHandler 消息管理抽象层
type IMsgHandler interface {
	DoMsgHandler(request IRequest)          // 调度/执行对应的 Router 消息处理方法
	AddRouter(msgID uint32, router IRouter) // 为消息添加具体的处理逻辑
	StartWorkerPool()                       // 启动一个Worker工作池
	SendMsgToTaskQueue(request IRequest)    // 将消息发送给消息任务队列处理
}
