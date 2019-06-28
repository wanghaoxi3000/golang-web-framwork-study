package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// Server IServer 的接口实现，定义一个 Server 的服务器模块
type Server struct {
	Name       string             // 服务器的名称
	IPVersion  string             // 服务器绑定的ip版本
	IP         string             // 服务器监听的IP
	Port       int                // 服务器监听的端口
	MsgHandler ziface.IMsgHandler // 当前的 Server 添加一个 router，注册连接接对应的处理业务
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf("[init] Server Name : %s, listenner at IP : %s, Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TCPPort)
	fmt.Printf("[init] Version %s, MaxConn:%d, MaxPackeetSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {
		// 0 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1 获取一个 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error : ", err)
			return
		}

		// 2 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("start Zinx server success, ", s.Name, " succ, Listenning...")
		var cid uint32
		cid = 0

		// 3 阻塞的等待客户端连接，处理客户端连接业务(读写)
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 将处理新连接的业务方法 和 conn 进行绑定 得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

// Serve 运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

// Stop 停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息 进行停止或者回收
}

// AddRouter 给当前的服务注册一个路由方法，供客户端的连接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!!")
}

// NewServer 初始化Server模块的方法
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		MsgHandler: NewMsgHandle(),
	}

	return s
}
