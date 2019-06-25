package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/ziface"
)

// Server IServer 的接口实现，定义一个 Server 的服务器模块
type Server struct {
	Name      string // 服务器的名称
	IPVersion string // 服务器绑定的ip版本
	IP        string // 服务器监听的IP
	Port      int    // 服务器监听的端口
}

// CallBackToClient 定义当前客户端链接的所绑定handle api (以后优化应该由用户自定义handle方法)
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallbackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}

	return nil
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP :%s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
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
			dealConn := NewConnection(conn, cid, CallBackToClient)
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

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return s
}
