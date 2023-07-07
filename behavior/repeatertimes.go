package behavior

import (
	"errors"
	"github.com/gamelee/robot/di"
)

type RepeaterTimes struct {
	Decorator
	Times int
}

func (rt *RepeaterTimes) Title() string {
	return "循环"
}

func (rt *RepeaterTimes) OnStart(injector di.Injector) (interface{}, error) {
	if rt.GetChild() == nil {
		return nil, errors.New("没有子节点")
	}
	var (
		rtn interface{}
		err error
	)
	//次数循环
	for ; rt.Times != 0; rt.Times-- {
		rtn, err = rt.GetChild().Run(injector)
		if err != nil {
			return nil, err
		}
	}
	return rtn, nil
}
