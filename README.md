golang 简易游戏行为树
================

灵感来源于 b3go, 去除了 tick 逻辑，用于测试游戏服务端逻辑

配合行为树编辑器和客户端，可以快速开发游戏业务测试

# Example

```go
package main

import (
	"github.com/gamelee/robot/behavior"
	"github.com/gamelee/robot/di"
	"log"
	"strings"
)

const config = `{
	"e5da3877-f8b4-4954-8279-61fb224ca829": {
        "id": "e5da3877-f8b4-4954-8279-61fb224ca829",
        "title": "测试",
        "root": "fb7a0e9c-7d80-4701-938d-66a90c8ec8dc",
        "properties": {},
        "nodes": {
            "fb7a0e9c-7d80-4701-938d-66a90c8ec8dc": {
                "id": "fb7a0e9c-7d80-4701-938d-66a90c8ec8dc",
                "name": "Sequence",
                "category": "composite",
                "title": "根节点",
                "properties": {},
                "children": [
                    "e665244d-5d04-40f2-b4bd-7ab8c8d880d0"
                ]
            },
            "e665244d-5d04-40f2-b4bd-7ab8c8d880d0": {
                "id": "e665244d-5d04-40f2-b4bd-7ab8c8d880d0",
                "title": "发呆",
                "name": "Nop",
                "category": "action",
                "properties": {
                    "Server": "pvp"
                }
            }
        }
    }
}`

func main() {

	projectConfig, err := behavior.LoadProjectConfig(strings.NewReader(config))
	if err != nil {
		log.Fatalln("load", err)
	}
	proj := behavior.NewProject()
	err = proj.Init(projectConfig)
	if err != nil {
		log.Fatalln("init", err)
	}
	_, err = proj.GetTreeByName("测试").Run(di.NewInjector())
	if err != nil {
		log.Fatalln("run", err)
	}
}

```
