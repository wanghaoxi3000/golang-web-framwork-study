package ziface

//IMessage 将请求的数据封装到一个 Message 消息中
type IMessage interface {
	GetMsgID() uint32  //获取消息的ID
	GetMsgLen() uint32 //获取消息的长度
	GetData() []byte   //获取消息的内容
	SetMsgID(uint32)   //设置消息的ID
	SetData([]byte)    //设置消息的内容
	SetDataLen(uint32) //设置消息的长度
}
