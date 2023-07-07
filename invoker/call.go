// Copyright(C),2020-2025
// Author:  lijie
// Version: 1.0.0
// Date:    2021/8/27 17:26
// Description:

package invoker

import (
	"encoding/json"
	"github.com/gamelee/robot"
	"sync/atomic"
)

type (
	Call struct {
		Seq     int64       `json:"seq"`
		Action  string      `json:"action"`
		Message string      `json:"msg"`
		Arg     interface{} `json:"arg"`
	}
	CallRst struct {
		Seq     int64       `json:"seq"`
		Code    Code        `json:"code"`
		Message string      `json:"msg"`
		Data    interface{} `json:"data"`
	}
)

var seq int64

func NewCall(action string, message string, args ...interface{}) *Call {
	this := &Call{Action: action, Message: message, Arg: args, Seq: atomic.AddInt64(&seq, 1)}
	if len(args) == 1 {
		this.Arg = args[0]
	}
	return this
}
func (c *Call) String() string {
	str, _ := robot.JsonEncodeString(c)
	return str
}

func NewCallRst(data []byte) *CallRst {
	rst := new(CallRst)
	err := json.Unmarshal(data, rst)
	if err != nil {
		rst.Code = CodeFailed
		rst.Message = err.Error()
	}
	return rst
}
