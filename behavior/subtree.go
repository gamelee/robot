package behavior

import (
	"errors"
	"fmt"

	"github.com/gamelee/robot/di"
)

// SubTree 子树，通过Name关联树ID查找
type SubTree struct {
	Action
	Tree string
}

func (st *SubTree) Title() string {
	return "子树"
}

func (st *SubTree) OnStart(injector di.Injector) (interface{}, error) {
	if st.proj == nil {
		return nil, errors.New("未找到项目信息")
	}
	tree := st.proj.GetTree(st.Tree)
	if tree == nil {
		tree = st.proj.GetTreeByName(st.Tree)
		if tree != nil {
			return nil, fmt.Errorf("load tree: %s", st.Tree)
		}
	}
	return tree.Run(injector)
}

type SubTreeLoadFunc func(string) (*Tree, error)
