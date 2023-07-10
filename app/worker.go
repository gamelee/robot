package main

import (
	"github.com/gamelee/robot"
	"github.com/gamelee/robot/app/msg"
	"github.com/gamelee/robot/invoker"
	"runtime"
)

type DDChessWorker struct {
	*robot.BaseWorker
}

func (chess *DDChessWorker) JSReqAll() (map[string]interface{}, error) {
	return msg.GetAllReqBody(), nil
}

func (chess *DDChessWorker) JSReqSend(dat map[string]interface{}) (interface{}, error) {
	panic("implements me")
}

func NewDDChessWorker() *DDChessWorker {
	return &DDChessWorker{BaseWorker: robot.NewBaseWorker()}
}

func (chess *DDChessWorker) JSGolangVersion() ([]byte, error) {
	return []byte(runtime.Version()), nil
}

// msgReceive
// 将收到的消息通知前端 js
// 前端返回 true 表示已收到，完成整个消息的发送
// 部分 go节点中需要根据 js 中的状态来执行，此逻辑保证了时序的
func (chess *DDChessWorker) msgReceiver() {
	for message := range chess.Hub.EventChan() {
		chess.Notify2JS(invoker.NewCall("NotifyJsMsg", "收发包消息", message))
	}
}
