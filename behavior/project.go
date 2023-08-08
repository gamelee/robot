// Copyright(C),2020-2025
// Author:  lijie
// Version: 1.0.0
// Date:    2021/8/27 9:32
// Description:

package behavior

import (
	"sync"
)

type Project struct {
	file  string
	Conf  ProjectConfig
	trees sync.Map
}

func (proj *Project) Init(conf ProjectConfig) error {
	proj.Conf = conf
	// 实例化行为树
	for i := range proj.Conf {
		treeConf := proj.Conf[i]
		tree := NewTree(proj)
		err := tree.Load(treeConf, defaultRegMap)
		if err != nil {
			return err
		}
		proj.trees.Store(treeConf.ID, tree)
	}
	return nil
}

func (proj *Project) GetTree(idOrTitle string) *Tree {

	tree, ok := proj.trees.Load(idOrTitle)
	if ok {
		return tree.(*Tree)
	}
	var rTree *Tree
	proj.trees.Range(func(_, tree interface{}) bool {
		if tree.(*Tree).Conf.Title == idOrTitle {
			rTree = tree.(*Tree)
			return false
		}
		return true
	})
	return rTree
}

func (proj *Project) GetTreeByName(name string) *Tree {
	var rTree *Tree
	proj.trees.Range(func(_, value interface{}) bool {
		tree, ok := value.(*Tree)
		if ok {
			return true
		}
		if tree.Conf.Title == name {
			rTree = tree
			return false
		}
		return true
	})
	return rTree
}

func NewProject() *Project {
	p := &Project{
		trees: sync.Map{},
	}
	return p
}
