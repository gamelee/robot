package behavior

type IComposite interface {
	IBaseNode
	GetChildCount() int
	GetChild(index int) IBaseNode
	AddChild(child IBaseNode)
}

type Composite struct {
	Action
	BaseWorker
	children []IBaseNode
}

func (cmp *Composite) Cate() string {
	return COMPOSITE
}

func (cmp *Composite) Title() string {
	return "组合节点"
}

func (cmp *Composite) GetChildCount() int {
	return len(cmp.children)
}

func (cmp *Composite) GetChild(index int) IBaseNode {
	return cmp.children[index]
}

func (cmp *Composite) AddChild(child IBaseNode) {
	cmp.children = append(cmp.children, child)
}

var _ IBaseNode = (*Composite)(nil)
