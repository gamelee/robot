package connect

import (
	"github.com/gamelee/robot"
	"sync"
)

type MessageType string

const (
	ReqMessage       MessageType = `req`
	RspMessage       MessageType = `rsp`
	NtfMessage       MessageType = `ntf`
	TreeEnterMessage MessageType = `enter`
	TreeOutMessage   MessageType = `out`
	SysMessage       MessageType = `sys`
)

type Message struct {
	Type     MessageType `json:"type"`
	From     string      `json:"from"`
	ID       interface{} `json:"id,omitempty"`
	IDPretty string      `json:"id_pretty,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	WaitId   interface{} `json:"wait_id,omitempty"`
	Error    error       `json:"error,omitempty"`
}

func NewMessage(from string, id interface{}, data interface{}) *Message {
	return &Message{From: from, ID: id, Data: data}
}
func NewSysMessage(from string, id interface{}, data interface{}) *Message {
	return &Message{From: from, Type: SysMessage, ID: id, Data: data}
}

func (m *Message) String() string {
	str, _ := robot.JsonEncodeString(m)
	return str
}

type MessageReg struct {
	m *sync.Map
}

func NewMessageReg() *MessageReg {
	return &MessageReg{m: &sync.Map{}}
}

func (mr *MessageReg) Reg(id interface{}, ch chan *Message) {
	mr.m.Store(id, ch)
}

func (mr *MessageReg) Del(id interface{}) {
	mr.m.Delete(id)
}

func (mr *MessageReg) Get(id interface{}) chan *Message {
	ch, ok := mr.m.Load(id)
	if !ok {
		return nil
	}
	return ch.(chan *Message)
}

func (mr *MessageReg) GetOrReg(id interface{}, ch chan *Message) chan *Message {
	c, _ := mr.m.LoadOrStore(id, ch)
	return c.(chan *Message)
}

func (mr *MessageReg) Range(fn func(k, v interface{}) bool) {
	mr.m.Range(fn)
}
