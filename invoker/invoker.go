// Copyright(C),2020-2025
// Author:  lijie
// Version: 1.0.0
// Date:    2021/8/27 17:30
// Description:

package invoker

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

type FuncManager struct {
	fs   map[string]*Function
	lock sync.RWMutex
}

func NewFuncManager() *FuncManager {
	return &FuncManager{
		fs: make(map[string]*Function),
	}
}

// RegStruct 注册结构体中的方法
func (fm *FuncManager) RegStruct(handle interface{}, filter func(structName, methodName string) (name string, use bool)) error {
	rv := reflect.ValueOf(handle)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return errors.New("请传入结构体或结构体指针")
	}
	rv = reflect.ValueOf(handle)
	rt := reflect.TypeOf(handle)
	name := rt.String()
	for i := 0; i < rt.NumMethod(); i++ {
		rvm := rv.Method(i)
		rtm := rt.Method(i)

		regName, use := filter(name, rtm.Name)
		if !use {
			continue
		}
		fn := newFunction(rtm.Type, rvm).SetName(regName)
		fn.in = rtm.Type.NumIn() - 1
		fm.Reg(fn)
	}
	return nil
}

// Reg 注册方法
func (fm *FuncManager) Reg(fn *Function) {
	fm.lock.Lock()
	fm.fs[fn.Name()] = fn
	fm.lock.Unlock()
}

// Invoke 调用已注册的方法, 对 invoker 进行了封装
func (fm *FuncManager) Invoke(name string, cb interface{}, args ...interface{}) error {
	rst := fm.Call(NewCall(name, "call:"+name, args...))

	if rst.Code != CodeSuccess {
		return errors.New(rst.Message)
	}
	fn, err := NewFunc(cb)
	if err != nil {
		return err
	}
	_, err = fn.Call(rst.Data)
	return err
}

// Call 调用已注册的方法
func (fm *FuncManager) Call(c *Call) (rst *CallRst) {
	rst = new(CallRst)
	rst.Seq = c.Seq

	fm.lock.RLock()
	defer fm.lock.RUnlock()
	fn, has := fm.fs[c.Action]
	if !has {
		rst.Code = CodeFailed
		rst.Message = "未找到方法:" + c.Action
		return
	}
	defer func() {
		if r := recover(); r != nil {
			rst.Code = CodeFailed
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			rst.Message = fmt.Sprintf("panic:%+v, Stack:%s", r, string(buf[:n]))
		}
	}()
	rtn, err := fn.Call(c.Arg)
	if err != nil {
		rst.Code = CodeFailed
		rst.Message = err.Error()
	} else {
		rst.Code = CodeSuccess
		rst.Message = "调用成功"
	}
	if len(rtn) == 1 {
		rst.Data = rtn[0]
	} else {
		rst.Data = rtn
	}
	return
}
