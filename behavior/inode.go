package behavior

import (
	"errors"

	"github.com/gamelee/robot/di"
)

type IBaseRunner interface {
	Title() string
	Run(note di.Injector) (interface{}, error)
}

type IBaseNode interface {
	IBaseRunner

	Init(conf *NodeConfig, project *Project)
	Conf() *NodeConfig
	Cate() string
	NodeWorker(...IBaseWorker) IBaseWorker
}

type IBaseWorker interface {
	OnStart(note di.Injector) (interface{}, error)
}

type BaseWorker struct{}

func (bw *BaseWorker) OnStart(_ di.Injector) (interface{}, error) {
	return nil, errors.New("no BaseWorker")
}
