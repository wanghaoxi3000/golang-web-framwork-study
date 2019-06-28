package ziface

// IServer 定义一个服务器接口
type IServer interface {
	Start()                                 // 启动服务器
	Stop()                                  // 停止服务器
	Serve()                                 // 运行服务器
	AddRouter(msgID uint32, router IRouter) // 路由功能 给当前的服务注册一个路由方法，供客户端的连接处理使用
	GetConnMgr() IConnManager               // 获取连接管理器

	SetOnConnStart(func(conneciton IConnection)) // 注册OnConnStart 钩子函数的方法
	SetOnConnStop(func(conneciton IConnection))  // 注册OnConnStop钩子函数的方法
	CallOnConnStart(conneciton IConnection)      // 调用OnConnStart钩子函数的方法
	CallOnConnStop(conneciton IConnection)       // 调用OnConnStop钩子函数的方法
}
