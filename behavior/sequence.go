package behavior

import (
	"github.com/gamelee/robot/di"
)

type Sequence struct {
	Composite
}

func (s *Sequence) Title() string {
	return "顺序执行"
}

func (s *Sequence) OnStart(injector di.Injector) (interface{}, error) {
	rtn := make([]interface{}, s.GetChildCount())
	for i := 0; i < s.GetChildCount(); i++ {
		data, err := s.GetChild(i).Run(injector)
		rtn[i] = data
		if err != nil {
			return rtn, err
		}
	}
	return rtn, nil
}
