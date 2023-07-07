package behavior

import (
	"fmt"
	"reflect"
)

type Node struct {
	Name  string                 `json:"name,omitempty"`
	Title string                 `json:"title,omitempty"`
	Cate  string                 `json:"category,omitempty"`
	Prop  map[string]interface{} `json:"properties,omitempty"`
}

// RegisterStructMaps 定义注册结构map
type RegisterStructMaps struct {
	maps map[string]reflect.Type
	info map[string]*Node
}

func NewRegisterStructMaps() *RegisterStructMaps {
	return &RegisterStructMaps{
		maps: make(map[string]reflect.Type),
		info: make(map[string]*Node),
	}
}

// CreateElem 根据 name 初始化结构
// 在这里根据结构的成员注解进行DI注入，这里没有实现，只是简单都初始化
func (rsm *RegisterStructMaps) CreateElem(name string) (interface{}, error) {
	var c interface{}
	var err error
	if v, ok := rsm.maps[name]; ok {
		c = reflect.New(v).Interface()
		return c, nil
	} else {
		err = fmt.Errorf("not found %s struct", name)
	}
	return nil, err
}

// CheckElem 查询是否存在
func (rsm *RegisterStructMaps) CheckElem(name string) bool {
	if _, ok := rsm.maps[name]; ok {
		return true
	}
	return false
}

// Register 根据类名注册实例
func (rsm *RegisterStructMaps) Register(c IBaseNode) {
	v := reflect.TypeOf(c).Elem()
	k := v.Name()

	if rsm.maps[k] != nil || rsm.info[c.Title()] != nil {
		panic(v.Name() + "," + c.Title() + "is exists")
	}
	rsm.maps[k] = v
	node := &Node{
		Name:  v.Name(),
		Title: c.Title(),
		Cate:  c.Cate(),
		Prop:  make(map[string]interface{}),
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Tag.Get("inject") != "" {
			continue
		}
		n := f.Type.String()
		if n == "behavior.Action" || n == "behavior.Composite" || n == "behavior.Decorator" {
			continue
		}
		if f.Type.Kind() == reflect.Func || f.Type.Kind() == reflect.Chan {
			continue
		}
		node.Prop[f.Name] = reflect.Zero(f.Type).Interface()
	}
	rsm.info[k] = node
}

// Combine 根据类名注册实例
func (rsm *RegisterStructMaps) Combine(reg *RegisterStructMaps) *RegisterStructMaps {
	for name, elem := range reg.maps {
		rsm.maps[name] = elem
	}
	return rsm
}

func (rsm *RegisterStructMaps) Nodes() map[string]*Node {
	return rsm.info
}
