package behavior

import (
	"encoding/json"
	"fmt"
	"github.com/gamelee/robot/di"
)

type BaseNode struct {
	IBaseWorker
	conf  *NodeConfig
	proj  *Project
	Error error
	Data  interface{}
}

func (bn *BaseNode) NodeWorker(workers ...IBaseWorker) IBaseWorker {
	if len(workers) > 0 {
		bn.IBaseWorker = workers[0]
	}
	return bn.IBaseWorker
}

func (bn *BaseNode) Conf() *NodeConfig {
	return bn.conf
}

func (bn *BaseNode) Init(conf *NodeConfig, proj *Project) {
	bn.conf = conf
	bn.proj = proj
	buf, _ := json.Marshal(conf.Properties)
	_ = json.Unmarshal(buf, bn.NodeWorker())
}

func (bn *BaseNode) Run(injector di.Injector) (rtn interface{}, err error) {
	_, err = injector.Invoke(func(emit func(life NodeLife, conf *NodeConfig, msg string) bool) (interface{}, error) {
		emit(NodeLifeStart, bn.conf, "")
		defer func() {
			if err != nil {
				emit(NodeLifeError, bn.conf, err.Error())
			}
		}()
		err = injector.Apply(bn.NodeWorker())
		if err != nil {
			err = fmt.Errorf(`配置节点 "%s" 出错: %w`, bn.conf.Title, bn.Error)
			return nil, err
		}
		rtn, err = bn.OnStart(injector)
		emit(NodeLifeFinish, bn.conf, "")
		return rtn, err
	})
	return
}
