package behavior

type IDecorator interface {
	IBaseNode
	SetChild(child IBaseNode)
	GetChild() IBaseNode
}

type Decorator struct {
	Action
	BaseWorker
	child IBaseNode
}

func (d *Decorator) Cate() string {
	return DECORATOR
}

func (d *Decorator) Title() string {
	return "包装节点"
}

func (d *Decorator) GetChild() IBaseNode {
	return d.child
}

func (d *Decorator) SetChild(child IBaseNode) {
	d.child = child
}

var _ IBaseNode = (*Decorator)(nil)
