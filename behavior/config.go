package behavior

import (
	"encoding/json"
	"fmt"
	"io"
)

// ProjectConfig 原生工程json类型
type ProjectConfig map[string]*TreeConfig

// LoadProjectConfig 加载原生工程
func LoadProjectConfig(r io.Reader) (ProjectConfig, error) {

	var project ProjectConfig
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &project)
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
