// Copyright(C),2020-2025
// Author:  lijie
// Version: 1.0.0
// Date:    2021/9/3 17:15
// Description: 用户抽象

package robot

import (
	"encoding/json"
	"github.com/gamelee/robot/behavior"
	"github.com/gamelee/robot/connect"
	"github.com/gamelee/robot/di"
	"github.com/gamelee/robot/store"
)

type Robot struct {
	// 链接管理器
	*connect.Hub
	// 行为树管理器
	*behavior.Project
	// 数据存储
	Blk *store.Blackboard
	// 依赖注入器
	di.Injector

	Data     interface{}
	Err      error
	LastTree string
}

func (rob *Robot) MarshalJSON() ([]byte, error) {
	out := make(map[string]interface{})
	//for _, key := range []string{zizouqi.IdxDeviceID, zizouqi.IdxPlayerID, zizouqi.IdxUserName} {
	//	out[key] = rob.Blk.Get(key)
	//}
	out["data"] = rob.Data
	out["error"] = rob.Err
	return json.Marshal(out)
}

func (rob *Robot) String() string {
	buf, _ := rob.MarshalJSON()
	return string(buf)
}

// NewRobot 创建机器人
func NewRobot() *Robot {
	this := &Robot{
		Injector: di.NewInjector(),
		Hub:      connect.NewManager(),
		Blk:      store.NewBlackboard(),
		Project:  behavior.NewProject(),
	}
	this.Map(this.Blk)
	this.Map(this.Hub)
	return this
}

// Reload
// @Description: 重新加载行为树文件和配置信息
// @receiver *Robot
// @param file 行为树文件
// @param data 配置信息
// @return error
func (rob *Robot) Reload(conf behavior.ProjectConfig, data map[string]interface{}) error {
	err := rob.Init(conf)
	if err != nil {
		return err
	}
	for k, v := range data {
		rob.Blk.Set(k, v)
	}
	return nil
}

// Run
// @Description: 运行行为树
// @receiver *Robot
// @param treeID 行为树 ID
// @return interface{} 运行结果
// @return error
func (rob *Robot) Run(treeID string) error {
	rob.LastTree = treeID
	tree, err := rob.GetTree(rob.LastTree)
	if tree == nil {
		return err
	}
	rob.Data, rob.Err = tree.Run(rob.Injector)
	return rob.Err
}
