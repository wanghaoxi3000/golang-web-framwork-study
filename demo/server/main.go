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

// PreHandle 在处理 conn 业务之前的钩子方法 Hook
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Handle 在处理 conn 业务的主方法 hook
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Ping...ping...\n"))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

// PostHandle 在处理 conn 业务知乎的钩子方法 Hook
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	s := znet.NewServer("[zinx v0.3]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
