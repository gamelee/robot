package connect

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

var (
	ErrTimeOut   = errors.New("timeout")
	RepeatedWait = errors.New("repeated wait")
)

type wrap struct {
	Connector
	hds       *MessageReg
	ch        chan *Message
	wg        *sync.WaitGroup
	lastWrite time.Time
	heart     chan struct{}
}

func checkErrorBreak(err error) bool {

	if err == nil {
		return false
	}
	if err == io.EOF || err == net.ErrClosed {
		return true
	}
	if e, ok := err.(*net.OpError); ok && !e.Temporary() {
		return true
	}
	return false
}

func (w *wrap) reader(ch chan *Message) (err error) {

	defer func() {
		w.wg.Done()
		ch <- NewSysMessage(w.Name(), "exit", "reader")
	}()
	ch <- NewSysMessage(w.Name(), "run", "reader")
	for {
		message := w.Read()
		if checkErrorBreak(message.Error) {
			err = message.Error
			return
		}

		ch <- message
		if message.WaitId != nil && message.WaitId != 0 {
			if c := w.hds.Get(message.WaitId); c != nil {
				c <- message
				w.hds.Del(message.WaitId)
				close(c)
			}
		}
		if !w.Running() {
			return
		}
	}
}

func (w *wrap) writer(ch chan *Message) (err error) {

	defer func() {
		w.wg.Done()
		ch <- NewSysMessage(w.Name(), "exit", "writer")
	}()

	ch <- NewSysMessage(w.Name(), "run", "writer")

	for message := range w.ch {

		message.Error = w.Write(message)

		if checkErrorBreak(message.Error) {
			err = message.Error
			return
		}
		ch <- message
		if !w.Running() {
			return
		}
		w.lastWrite = time.Now()
	}
	return nil
}

func (w *wrap) heartbeat(ch chan *Message) (err error) {
	ka, ok := w.Connector.(KeepAlive)
	if !ok || ka == nil {
		return nil
	}
	defer func() {
		ch <- NewSysMessage(w.Name(), "exit", "heart")
	}()
	ch <- NewSysMessage(w.Name(), "run", "heart")
	delta, message := ka.KeepAlive()
	tick := time.NewTicker(delta)
	for {
		select {
		case <-w.heart:
			return
		case <-tick.C:
		}
		if !w.Running() {
			return
		}
		//sub := time.Now().Sub(w.lastWrite)
		//if  {
		//	tick.Reset(delta - sub)
		//}
		_, message = ka.KeepAlive()

		if checkErrorBreak(message.Error) {
			return
		}
		w.ch <- message
	}
}

func (w *wrap) Stop() {
	if w.Running() {
		w.Connector.Close()
		close(w.ch)
		close(w.heart)
		w.wg.Wait()
	}
	log.Println("stop", w.Name())
}

// Hub 连接管理器
type Hub struct {
	lock    sync.RWMutex
	clients map[string]*wrap
	ch      chan *Message
}

func NewManager() *Hub {
	return &Hub{
		lock:    sync.RWMutex{},
		clients: make(map[string]*wrap),
		ch:      make(chan *Message, 100),
	}
}

func (h *Hub) Add(conn Connector) error {

	if h.get(conn.Name()) != nil {
		panic("服务已存在:" + conn.Name())
	}
	w := &wrap{
		ch:        make(chan *Message, 12),
		hds:       NewMessageReg(),
		Connector: conn,
		wg:        &sync.WaitGroup{},
		heart:     make(chan struct{}),
	}
	h.lock.Lock()
	h.clients[w.Name()] = w
	h.lock.Unlock()
	return h.startClient(w)
}

func (h *Hub) get(cliName string) *wrap {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.clients[cliName]
}

func (h *Hub) Stop(cliName, reason string) {
	if h == nil {
		return
	}
	w := h.get(cliName)
	delete(h.clients, cliName)
	if w == nil {
		return
	}
	w.Stop()
}

func (h *Hub) StopAll(reason string) {
	if h == nil {
		return
	}
	for srvName := range h.clients {
		h.Stop(srvName, reason)
	}
}

func (h *Hub) ReqWait(cliName string, message *Message, timeout ...time.Duration) *Message {
	w := h.get(cliName)
	chMessage := w.wait(message.WaitId)
	w.ch <- message
	defer func() { w.hds.Del(message.WaitId) }()
	if message.Error != nil {
		return message
	}
	if len(timeout) == 0 {
		return <-chMessage
	}
	select {
	case rst := <-chMessage:
		return rst
	case <-time.After(timeout[0]):
		message.Error = fmt.Errorf("reqwait message %v time out %s", message.WaitId, timeout[0])
		return message
	}
}

func (h *Hub) NewMessage(cliName string, id interface{}, req proto.Message) *Message {
	w := h.get(cliName)
	if w == nil {
		m := NewSysMessage(cliName, id, req)
		m.Error = fmt.Errorf("%s 未建立链接", cliName)
		return m
	}
	return w.NewMessage(id, req)
}

// Req 发送消息
func (h *Hub) Req(cliName string, id interface{}, req proto.Message) {
	w := h.get(cliName)
	w.ch <- w.NewMessage(id, req)
}

func (h *Hub) Wait(cliName string, waitId interface{}, timeout ...int64) *Message {
	w := h.get(cliName)
	chMessage := w.wait(waitId)
	if len(timeout) == 0 {
		return <-chMessage
	}
	select {
	case rst := <-chMessage:
		return rst
	case <-time.After(time.Duration(timeout[0]) * time.Second):
		message := w.NewMessage(waitId, ErrTimeOut)
		message.Error = fmt.Errorf("wait message %v time out %d", message.WaitId, timeout)
		return message
	}
}

// wait 等待指定id的消息
func (w *wrap) wait(waitId interface{}) chan *Message {

	var (
		ch = make(chan *Message, 1)
	)

	if w.hds.Get(waitId) != nil {
		ch <- w.NewMessage(waitId, RepeatedWait)
		return ch
	}

	w.hds.Reg(waitId, ch)
	return ch
}

// startClient
// @Description: 建立新的链接
// @receiver this
// @param wrap
// @return err
func (h *Hub) startClient(w *wrap) (err error) {
	err = w.Connect()
	if err != nil {
		return err
	}
	w.wg.Add(2)

	// 读消息循环
	go func() {
		recovery(func() error {
			return w.reader(h.ch)
		})
	}()
	go func() {
		recovery(func() error {
			return w.writer(h.ch)
		})
	}()
	go func() {
		recovery(func() error {
			return w.heartbeat(h.ch)
		})
	}()

	log.Printf("start client %v", w)
	return nil
}

func (h *Hub) EventChan() chan *Message {
	return h.ch
}

func recovery(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
		if err != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("error: %w;\ntrace: %v", err, bytes2String(buf))
		}
	}()
	err = f()
	return
}
