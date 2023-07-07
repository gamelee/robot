package behavior

import (
	"encoding/json"
	"fmt"
	"github.com/gamelee/robot"
	"os"
)

// ProjectConfig 原生工程json类型
type ProjectConfig map[string]*TreeConfig

// LoadProjectConfig 加载原生工程
func LoadProjectConfig(path string) (ProjectConfig, error) {

	var project ProjectConfig
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &project)
	if err != nil {
		fmt.Println("parse failed", err)
		return nil, err
	}

	return project, nil
}

type TreeConfig struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Root       string                 `json:"Root"`
	Properties map[string]interface{} `json:"properties"`
	Nodes      map[string]NodeConfig  `json:"nodes"`
}

// NodeConfig
// 编辑器地址
// 节点json类型
type NodeConfig struct {
	Id         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Category   string                 `json:"category,omitempty"`
	Title      string                 `json:"title,omitempty"`
	Children   []string               `json:"children,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

func (nc *NodeConfig) String() string {
	buf, _ := json.Marshal(nc)
	return robot.Bytes2String(buf)
}
