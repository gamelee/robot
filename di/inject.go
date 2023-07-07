package di

import (
	"fmt"
	"reflect"
)

type Injector interface {
	Applicator
	Invoker
	TypeMapper
	SetParent(Injector)
}

type Applicator interface {
	Apply(interface{}) error
}

type Invoker interface {
	Invoke(interface{}) ([]reflect.Value, error)
}

type TypeMapper interface {
	Map(interface{}) TypeMapper
	MapTo(interface{}, interface{}) TypeMapper
	Set(reflect.Type, reflect.Value) TypeMapper
	Get(reflect.Type) reflect.Value
}

func InterfaceOf(value interface{}) reflect.Type {
	rt := reflect.TypeOf(value)
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Interface {
		panic("Called inject.InterfaceOf with a value that is not a pointer to an interface")
	}
	return rt
}

type injector struct {
	values map[reflect.Type]reflect.Value
	parent Injector
}

func NewInjector() Injector {
	i := new(injector)
	i.values = make(map[reflect.Type]reflect.Value)
	return i
}

func (injector *injector) Invoke(f interface{}) ([]reflect.Value, error) {
	rt := reflect.TypeOf(f)
	args := make([]reflect.Value, rt.NumIn())
	for i := 0; i < rt.NumIn(); i++ {
		argType := rt.In(i)
		val := injector.Get(argType)
		if !val.IsValid() {
			return nil, fmt.Errorf("value not found for type %v", argType)
		}
		args[i] = val
	}
	return reflect.ValueOf(f).Call(args), nil
}

func (injector *injector) Apply(val interface{}) error {
	rv := reflect.ValueOf(val)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil // Should not panic here ?
	}
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)
		if fv.CanSet() && ft.Tag.Get("inject") == "true" {
			_ft := fv.Type()
			_v := injector.Get(_ft)
			if !_v.IsValid() {
				return fmt.Errorf("value not found for type %v", ft)
			}
			fv.Set(_v)
		}
	}
	return nil
}

func (injector *injector) SetParent(p Injector) {
	injector.parent = p
}

func (injector *injector) Map(val interface{}) TypeMapper {
	injector.values[reflect.TypeOf(val)] = reflect.ValueOf(val)
	return injector
}

func (injector *injector) MapTo(val interface{}, ifacePtr interface{}) TypeMapper {
	injector.values[InterfaceOf(ifacePtr)] = reflect.ValueOf(val)
	return injector
}

func (injector *injector) Set(typ reflect.Type, val reflect.Value) TypeMapper {
	injector.values[typ] = val
	return injector
}

func (injector *injector) Get(rt reflect.Type) reflect.Value {
	val := injector.values[rt]
	if val.IsValid() {
		return val
	}
	if rt.Kind() == reflect.Interface {
		for k, v := range injector.values {
			if k.Implements(rt) {
				val = v
				break
			}
		}
	}

	if !val.IsValid() && injector.parent != nil {
		val = injector.parent.Get(rt)
	}
	return val
}
