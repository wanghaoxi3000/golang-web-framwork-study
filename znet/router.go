package znet

import "zinx/ziface"

// BaseRouter 实现 router 时， 先嵌入这个 BaseRouter 基类， 然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

/*
 这里 BaseRouter 的方法都为空
 业务可根据需要实现对应的方法
*/

// PreHandle 在处理 conn 业务之前的钩子方法 Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// Handle 在处理 conn 业务的主方法 hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// PostHandle 在处理 conn 业务知乎的钩子方法 Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
