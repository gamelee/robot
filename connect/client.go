package connect

import "time"

type KeepAlive interface {
	KeepAlive() (time.Duration, *Message)
}


type Connector interface {
	Name() string
	Connect() error
	Close()
	Write(message *Message) error
	Read() *Message
	Running() bool
	NewMessage(id, data interface{}) *Message
}
