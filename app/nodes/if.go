// Copyright(C),2020-2025
// Author:  lijie
// Version: 1.0.0
// Date:    2021/8/26 18:35
// Description:

package nodes

import (
	"errors"
	"github.com/gamelee/robot/behavior"
	"github.com/gamelee/robot/connect"
	"github.com/gamelee/robot/di"
	"github.com/gamelee/robot/ui/lorca"
)

type IF struct {
	behavior.Composite
	EvalJs func(code string) lorca.Value `inject:"true"`
	Hub    *connect.Hub                  `inject:"true"`
	Cond   string
}

func (i *IF) Title() string {
	return "条件"
}
func (i *IF) OnStart(injector di.Injector) (interface{}, error) {
	if i.GetChildCount() != 2 {
		return nil, errors.New("子节点数量异常, except 2")
	}
	if i.Check() {
		return i.GetChild(0).Run(injector)
	} else {
		return i.GetChild(1).Run(injector)
	}
}

func (i *IF) Check() bool {

	if i.Cond == "true" || i.Cond == "1" {
		return true
	}

	r := i.EvalJs(i.Cond)
	return r.Bool()
}
