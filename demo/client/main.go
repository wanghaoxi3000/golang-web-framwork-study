package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)

	// 连接服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit")
		return
	}

	for {
		//发送封包 message 消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("ZinxV0.5 client test message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write conn err", err)
			return
		}

		//接收服务器回复的一个message数据
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error")
			break
		}

		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpacke err ", err)
			return
		}

		//根据head中的datalen 再读取data内容
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			//根据datalen的长度再次从io流中读取
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err: ", err)
				return
			}

			fmt.Println("---> Recv MsgID: ", msg.ID, ", datalen = ", msg.DataLen, "data = ", string(msg.Data))
		}

		time.Sleep(3 * time.Second)
	}

}
