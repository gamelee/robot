package main

import (
	"github.com/gamelee/robot/behavior"
	"log"
	"strings"
)

const config = `{
	"e5da3877-f8b4-4954-8279-61fb224ca829": {
        "id": "e5da3877-f8b4-4954-8279-61fb224ca829",
        "title": "断开pvp",
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
		panic(err)
		return
	}
	log.Println(projectConfig)
}
