package behavior

import (
	"github.com/gamelee/robot/di"
)

type Nop struct {
	Action
	Id  interface{}
	Id2 interface{}
}

func (nop *Nop) Title() string {
	return "空节点"
}

func (nop *Nop) OnStart(injector di.Injector) (interface{}, error) {
	return "nop", nil
}
