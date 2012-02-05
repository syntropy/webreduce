package messaging

import (
	zmq "github.com/alecthomas/gozmq"
)

var ctx zmq.Context

func init() {
	var err error
	ctx, err = zmq.NewContext()
	if err != nil {
		panic(err)
	}
}

type Socket interface {
	Bind(addr string) (err error)
	Connect(addr string) (err error)
	Close() (err error)
}

type Message struct {
	Payload []byte
}

func NewMessage(p []byte) (m *Message) {
	return &Message{Payload: p}
}

func NewStringMessage(body string) (m *Message) {
	return NewMessage([]byte(body))
}
