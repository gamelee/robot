package behavior

type IAction interface {
	IBaseNode
}

type Action struct {
	BaseNode
	BaseWorker
}

func (act *Action) Title() string {
	return "叶子节点"
}

func (act *Action) Cate() string {
	return ACTION
}

func (act *Action) Init(conf *NodeConfig, proj *Project) {
	act.BaseNode.Init(conf, proj)
}
