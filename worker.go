package robot

import (
	"errors"
	"github.com/gamelee/robot/behavior"
	"github.com/gamelee/robot/connect"
	"github.com/gamelee/robot/di"
	"github.com/gamelee/robot/invoker"
	"github.com/gamelee/robot/store"
	"os"
	"unsafe"
)

type NotifyJS = string

const (
	NotifyJsMsg NotifyJS = "msg"
)

type Worker interface {
	Notify2JS(call *invoker.Call) *invoker.CallRst

	setGame(game *Game)
	// JSReqAll 返回所有请求
	JSReqAll() (map[string]interface{}, error)
	// JSReqSend 发送请求
	JSReqSend(dat map[string]interface{}) (interface{}, error)
	// JSRobotRun 运行行为树
	JSRobotRun(treeFile, treeID string, data map[string]interface{}) (interface{}, error)
	// JSNodes 所有的 golang 节点
	JSNodes() (interface{}, error)
	// JSFileRead 读取文件
	JSFileRead(file string) (interface{}, error)
	// JSFileWrite 写入文件
	JSFileWrite(fn, text string) (interface{}, error)
}

type BaseWorker struct {
	game *Game
	// 链接管理器
	*connect.Hub
	// 行为树管理器
	*behavior.Project
	// 数据存储
	Blk *store.Blackboard
	// 依赖注入器
	di.Injector
}

func (b *BaseWorker) setGame(game *Game) {
	b.game = game
	b.Injector.Map(b.game.ui.RunJS)
}

func NewBaseWorker() *BaseWorker {
	w := &BaseWorker{
		Injector: di.NewInjector(),
		Hub:      connect.NewManager(),
		Blk:      store.NewBlackboard(),
		Project:  behavior.NewProject(),
	}
	//
	w.Map(w.Blk)
	w.Map(w.Hub)
	return w
}

func (b *BaseWorker) Notify2JS(call *invoker.Call) *invoker.CallRst {
	return b.game.ui.CallJS(call)
}

func (b *BaseWorker) JSRobotRun(treeFile, treeID string, data map[string]interface{}) (interface{}, error) {
	fd, err := os.Open(treeFile)
	if err != nil {
		return nil, err
	}
	conf, err := behavior.LoadProjectConfig(fd)
	if err != nil {
		return nil, err
	}
	err = b.Init(conf)
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		b.Blk.Set(k, v)
	}
	return b.GetTree(treeID).Run(b.Injector)
}

func (b *BaseWorker) JSNodes() (interface{}, error) {
	return behavior.Nodes(), nil
}

func (b *BaseWorker) JSFileRead(file string) (interface{}, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.New("读取文件出错:" + err.Error())
	}
	return bytes2String(buf), nil
}

func bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (b *BaseWorker) JSFileWrite(fn, text string) (interface{}, error) {
	err := os.WriteFile(fn, string2Bytes(text), os.ModePerm|os.ModeTemporary)
	if err != nil {
		return false, errors.New("写入文件出错:" + err.Error())
	}
	return true, nil
}

func string2Bytes(s string) []byte {
	tmp := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{tmp[0], tmp[1], tmp[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
