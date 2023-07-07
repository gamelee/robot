package behavior

import (
	"errors"
	"github.com/gamelee/robot/di"
	"time"
)

type Wait struct {
	Decorator
	During time.Duration
}

func (w *Wait) Title() string {
	return "等待"
}
func (w *Wait) OnStart(injector di.Injector) (interface{}, error) {

	if w.During <= 0 {
		return nil, errors.New("wait 时间不能小于 0")
	}
	select {
	case <-time.After(w.During * time.Millisecond):
	}
	if w.GetChild() == nil {
		return "没有子节点", nil
	}
	return w.GetChild().Run(injector)
}
