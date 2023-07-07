package store

import (
	"encoding/json"
	"sync"
)

type Blackboard struct {
	data *sync.Map
}

func (bb *Blackboard) MarshalJSON() ([]byte, error) {
	out := make(map[string]interface{})
	bb.data.Range(func(key, value interface{}) bool {
		out[key.(string)] = value
		return true
	})
	return json.Marshal(out)
}

func NewBlackboard() *Blackboard {
	return &Blackboard{data: &sync.Map{}}
}

func (bb *Blackboard) Get(name string) interface{} {
	t, _ := bb.data.Load(name)
	return t
}

func (bb *Blackboard) Has(name string) bool {
	_, ok := bb.data.Load(name)
	return ok
}

func (bb *Blackboard) Set(name string, val interface{}) *Blackboard {
	bb.data.Store(name, val)
	return bb
}

func (bb *Blackboard) String(args ...string) string {
	var name, dft string
	switch len(args) {
	case 0:
		return ""
	case 1:
		name = args[0]
	default:
		name, dft = args[0], args[1]
	}
	t := bb.Get(name)
	if t != nil {
		return t.(string)
	}
	return dft
}

func (bb *Blackboard) Int(name string) int {
	t := bb.Get(name)
	if t != nil {
		if it, ok := t.(float64); ok {
			return int(it)
		}
		return t.(int)
	}
	return 0
}

func (bb *Blackboard) Uint(name string) uint {
	t := bb.Get(name)
	if t != nil {
		if it, ok := t.(float64); ok {
			return uint(it)
		}
		return t.(uint)
	}
	return 0
}

func (bb *Blackboard) Bool(name string) bool {
	t := bb.Get(name)
	if t != nil {
		return t.(bool)
	}
	return false
}

func (bb *Blackboard) Int64(name string) int64 {
	t := bb.Get(name)
	if t != nil {
		if it, ok := t.(float64); ok {
			return int64(it)
		}
		return t.(int64)
	}
	return 0
}

func (bb *Blackboard) Uint32(name string) uint32 {
	t := bb.Get(name)
	if t != nil {
		if it, ok := t.(float64); ok {
			return uint32(it)
		}
		return t.(uint32)
	}
	return 0
}

func (bb *Blackboard) Uint64(name string) uint64 {
	t := bb.Get(name)
	if t != nil {
		if it, ok := t.(float64); ok {
			return uint64(it)
		}
		return t.(uint64)
	}
	return 0
}

func (bb *Blackboard) Del(name string) {
	bb.data.Delete(name)
}
func (bb *Blackboard) Bytes(name string) []byte {
	t := bb.Get(name)
	if t != nil {
		return t.([]byte)
	}
	return nil
}

func (bb *Blackboard) MapString(name string) map[string]string {
	t := bb.Get(name)
	if t != nil {
		return t.(map[string]string)
	}
	return nil
}

func (bb *Blackboard) MapInterface(name string) map[string]interface{} {
	t := bb.Get(name)
	if t != nil {
		return t.(map[string]interface{})
	}
	return nil
}

func (bb *Blackboard) Copy() *Blackboard {
	r := NewBlackboard()
	bb.data.Range(func(key, value interface{}) bool {
		r.Set(key.(string), value)
		return true
	})
	return r
}
