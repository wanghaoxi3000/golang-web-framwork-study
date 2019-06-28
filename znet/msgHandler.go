package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

//MsgHandle 消息处理模块的实现
type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter //存放每个MsgID 所对应的处理方法
	TaskQueue      []chan ziface.IRequest    //负责worker取任务的消息队列
	WorkerPoolSize uint32                    //Worker池的数量
}

//NewMsgHandle 创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前 msg 绑定的 API 处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api , msgID = " + strconv.Itoa(int(msgID)))
	}

	//添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " succ!")
}

//DoMsgHandler 执行对应的 Router 消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is NOT FOUND! Need Register!")
	}

	//2 根据MsgID 调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//StartWorkerPool 启动一个Worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//根据 workerPoolSize 分别开启 Worker，每个 Worker 用一个 go 来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1 当前的 worker 对应的 channel 消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前的 Worker， 阻塞等待消息从 channel 传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//StartOneWorker 启动一个 Worker 工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started ...")

	//不断的阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的Request, 执行当前Request所绑定业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//SendMsgToTaskQueue 将消息交给 TaskQueue， 由 Worker 进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//1 将消息平均分配给不通过的worker
	//根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" reqeust MsgID = ", request.GetMsgID(),
		" to WorkerID = ", workerID)

	//2 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
