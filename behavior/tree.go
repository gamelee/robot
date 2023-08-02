package behavior

import (
	"errors"
	"fmt"
	"github.com/gamelee/robot/di"
)

type NodeLife string

const (
	NodeLifeCreated NodeLife = `created`
	NodeLifeInit    NodeLife = `init`
	NodeLifeStart   NodeLife = `start`
	NodeLifeFinish  NodeLife = `finish  `
	NodeLifeError   NodeLife = `error  `
)

type NodeEvent struct {
	Action NodeLife
	Conf   *NodeConfig
	Msg    string
}

type Tree struct {
	IBaseRunner
	nodeEvent chan *NodeEvent
	Root      IBaseNode
	Conf      *TreeConfig
	proj      *Project
}

func NewTree(proj *Project) *Tree {
	tree := &Tree{
		nodeEvent: make(chan *NodeEvent, 256),
		proj:      proj,
	}
	return tree
}

func (tree *Tree) GetRoot() IBaseNode {
	return tree.Root
}

func (tree *Tree) Emit(life NodeLife, node *NodeConfig, msg string) bool {
	select {
	case tree.nodeEvent <- &NodeEvent{Action: life, Conf: node, Msg: msg}:
		return true
	default:
		return false
	}
}

func (tree *Tree) Event() chan *NodeEvent {
	return tree.nodeEvent
}

func (tree *Tree) Load(data *TreeConfig, maps *RegisterStructMaps) error {
	tree.Conf = data
	nodes := make(map[string]IBaseNode)
	for id := range tree.Conf.Nodes {
		conf := tree.Conf.Nodes[id]
		iNode, err := maps.CreateElem(conf.Name)
		if err != nil {
			return fmt.Errorf("未注册的节点: %+v", conf.Name)
		}
		node := iNode.(IBaseNode)
		node.NodeWorker(node.(IBaseWorker))
		tree.Emit(NodeLifeCreated, &conf, "")
		node.Init(&conf, tree.proj)
		nodes[id] = node
		tree.Emit(NodeLifeInit, &conf, "")
	}

	// 连接节点
	for id, nodeConf := range tree.Conf.Nodes {
		node := nodes[id]

		if node.Cate() == COMPOSITE && len(nodeConf.Children) > 0 {
			for i := 0; i < len(nodeConf.Children); i++ {
				var cid = nodeConf.Children[i]
				comp := node.(IComposite)
				comp.AddChild(nodes[cid])
			}
		} else if node.Cate() == DECORATOR && len(nodeConf.Children) > 0 {
			dec := node.(IDecorator)
			dec.SetChild(nodes[nodeConf.Children[0]])
		}
	}

	tree.Root = nodes[data.Root]
	return nil
}

func (tree *Tree) Run(injector di.Injector) (rst interface{}, err error) {
	if tree == nil {
		return nil, errors.New("tree is nil")
	}
	injector.Map(tree.Emit)
	return tree.Root.Run(injector)
}

var defaultRegMap = NewRegisterStructMaps()
