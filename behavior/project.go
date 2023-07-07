// Copyright(C),2020-2025
// Author:  lijie
// Version: 1.0.0
// Date:    2021/8/27 9:32
// Description:

package behavior

import (
	"errors"
	"sync"
)

type Project struct {
	file  string
	Conf  ProjectConfig
	trees sync.Map
}

func (proj *Project) init() error {
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

func (proj *Project) GetTree(id string) (*Tree, error) {

	tree, ok := proj.trees.Load(id)
	if ok {
		return tree.(*Tree), nil
	}
	var t *Tree
	proj.trees.Range(func(_, tree interface{}) bool {
		if tree.(*Tree).Conf.ID == id {
			t = tree.(*Tree)
			return false
		}
		return true
	})
	if t != nil {
		return t, nil
	}
	return nil, errors.New("未找到树 " + id)
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

func NewProject(file string) (*Project, error) {
	conf, err := LoadProjectConfig(file)
	if err != nil {
		return nil, err
	}
	p := &Project{
		file:  file,
		Conf:  conf,
		trees: sync.Map{},
	}
	return p, p.init()
}
