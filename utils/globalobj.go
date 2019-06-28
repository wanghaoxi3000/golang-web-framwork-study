package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*GlobalObj 存储一切可配置的全局参数，供其他模块使用
一些参数是可以通过 zinx.json 由用户进行配置
*/
type GlobalObj struct {
	// Server
	TCPServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TCPPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称

	// 框架信息
	Version          string // 当前版本号
	MaxConn          int    // 当前服务器主机允许的最大链接数
	MaxPackageSize   uint32 // 当前框架数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作Worker池Gorountine数量
	MaxWorkerTaskLen uint32 // 每个worker对应的消息队列的任务的数量最大值
}

//GlobalObject 定义一个全局的对外 Globalobj
var GlobalObject *GlobalObj

//Reload 从 zinx.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		// panic(err)
		return
	}
	//将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

//提供一个init方法，初始化当前的 GlobalObject
func init() {
	//如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.8",
		TCPPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	//尝试从 conf/zinx.json 去加载一些用户自定义的参数
	GlobalObject.Reload()
}
