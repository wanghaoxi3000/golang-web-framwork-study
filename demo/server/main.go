package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// PingRouter 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle 在处理 conn 业务的主方法 hook
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")

	//先读取客户端的数据，回写 ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	request.GetConnection().SendMsg(1, []byte("ping...ping...ping..."))
}

// HelloRouter 自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

// Handle 在处理 conn 业务的主方法 hook
func (p *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")

	//先读取客户端的数据，回写 ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	request.GetConnection().SendMsg(201, []byte("Hello..."))
}

//DoConnectBegin 创建连接之后执行的钩子函数
func DoConnectBegin(conn ziface.IConnection) {
	fmt.Println("===> DoConnectBegin is Called...")
	if err := conn.SendMsg(202, []byte("DoConnectBegin")); err != nil {
		fmt.Println(err)
	}
}

//DoConnectLost 连接断开之前执行的钩子函数
func DoConnectLost(conn ziface.IConnection) {
	fmt.Println("===> DoConnectLost is Called...")
	fmt.Println("conn ID = ", conn.GetConnID(), " is lost")
}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(DoConnectBegin)
	s.SetOnConnStop(DoConnectLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
