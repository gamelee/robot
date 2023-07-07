// Copyright(C),2020-2025
// Author:  lijie
// Version: 1.0.0
// Date:    2021/8/27 17:40
// Description:

package invoker

import (
	"fmt"
	"strconv"
	"testing"
)

func Demo1(a, b int) int {
	return a + b
}

func Demo2(c int) string {
	return strconv.Itoa(c)
}

type Demo struct {
	SN string
}

func Demo3(sn string) *Demo {
	return &Demo{sn}
}

func TestNewInvoker(t *testing.T) {
	i := NewFuncManager()
	fn, _ := NewFunc(Demo2)
	i.Reg(fn)
	err := i.Invoke("Demo2", Demo3, 3)
	if err != nil {
		return
	}
	t.Log()
}

func TestFunctionCall(t *testing.T) {
	fn, _ := NewFunc(Demo1)
	rst, err := fn.Call([]interface{}{1, 2})
	if err != nil {
		t.Error("call failed:", err)
		return
	}
	fmt.Println(rst)
}
