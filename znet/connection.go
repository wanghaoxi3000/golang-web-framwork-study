package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// Connection 连接模块
type Connection struct {
	Conn       *net.TCPConn       // 当前连接的 socket TCP 套接字
	ConnID     uint32             // 连接的 ID
	isClosed   bool               // 当前连接状态
	MsgHandler ziface.IMsgHandler // 当前连接所绑定的处理业务方法API
	ExitChan   chan bool          // 告知当前连接已经退出的/停止 channel
	msgChan    chan []byte        // 无缓冲管道，用于读、写 Goroutine 之间的消息通信
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}

	return c
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("[Reader is exit!], connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head 二级制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		//拆包，得到msgID 和 msgDatalen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		//根据dataLen  再次读取Data， 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前 conn 数据的 Request 请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给 Worker 工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 执行注册的路由方法，找到注册绑定的 Conn 对应的 router 调用
			// 根据绑定好的 MsgID 找到对应处理 api 业务执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

//StartWriter Goroutine, 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())

	//不断的阻塞的等待 channel 的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时 Writer 也要推出
			return
		}
	}
}

// Start 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID = ", c.ConnID)
	// 启动从当前连接的读数据的业务
	go c.StartReader()
	//启动从当前链接写数据的业务
	go c.StartWriter()
}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop... ConnID = ", c.ConnID)

	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭socket连接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// GetTCPConnection 获取当前连接绑定的 socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的 TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 发送数据， 将数据发送给远程的客户端
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	//将data进行封包 MsgDataLen|MsgID|Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}
