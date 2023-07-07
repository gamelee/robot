// Copyright(C),2020-2025
// Author:  lijie
// Version: 1.0.0
// Date:    2021/8/27 17:26
// Description:

package invoker

import (
	"encoding/json"
	"sync/atomic"
	"unsafe"
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
	str, _ := jsonEncodeString(c)
	return str
}

func jsonEncode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func jsonEncodeString(v interface{}) (string, error) {
	b, err := jsonEncode(v)
	if err != nil {
		return "", err
	}
	return bytes2String(b), nil
}

func bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
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
