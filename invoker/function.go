// Copyright(C),2020-2025
// Author:  lijie
// Version: 1.0.0
// Date:    2021/8/27 17:25
// Description:

package invoker

import (
	"errors"
	"reflect"
	"strconv"
)

type Code int

const (
	_           Code = iota
	CodeSuccess      = iota
	CodeFailed       = iota
)

type Function struct {
	name    string
	oriName string
	in      int
	out     int
	typ     reflect.Type
	val     reflect.Value
}

func (fn *Function) Name() string {
	return fn.name
}

func (fn *Function) SetName(name string) *Function {
	fn.name = name
	return fn
}

func newFunction(typ reflect.Type, val reflect.Value) *Function {
	return &Function{
		oriName: val.String(),
		in:      typ.NumIn(),
		out:     typ.NumOut(),
		typ:     typ,
		val:     val,
	}
}

func NewFunc(fn interface{}) (*Function, error) {
	t, v := reflect.TypeOf(fn), reflect.ValueOf(fn)
	if t.Kind() != reflect.Func {
		return nil, errors.New("必须注册函数类型")
	}
	return newFunction(t, v), nil
}

// Call 调用方法实现
func (fn *Function) Call(arg interface{}) (rtn []interface{}, error error) {
	args := make([]reflect.Value, fn.in)

	if fn.in == 1 {
		args[0] = reflect.ValueOf(arg)
	} else if len(args) > 1 {
		args2, ok := arg.([]interface{})
		if ok && len(args2) == len(args) {
			for i := 0; i < fn.in; i++ {
				args[i] = reflect.ValueOf(args2[i])
			}
		} else {
			return nil, errors.New("参数个数异常，期待个数:" + strconv.Itoa(len(args)))
		}
	}
	rst := fn.val.Call(args)
	for _, v := range rst {
		rtn = append(rtn, v.Interface())
	}
	return rtn, nil
}
