package ziface

// IServer 定义一个服务器接口
type IServer interface {
	Start()
	Stop()
	Serve()
}