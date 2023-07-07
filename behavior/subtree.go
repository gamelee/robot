package behavior

import (
	"errors"
	"fmt"

	"github.com/gamelee/robot/di"
)

// SubTree 子树，通过Name关联树ID查找
type SubTree struct {
	Action
	Load SubTreeLoadFunc `inject:"true"`
	Tree string
}

func (st *SubTree) Title() string {
	return "子树"
}

func (st *SubTree) OnStart(injector di.Injector) (interface{}, error) {

	if st.Load == nil {
		return nil, errors.New("未设置子树加载方法")
	}
	tree, err := st.Load(st.Tree)
	if err != nil {
		return nil, fmt.Errorf("load tree:%w", err)
	}
	return tree.Run(injector)
}

type SubTreeLoadFunc func(string) (*Tree, error)
