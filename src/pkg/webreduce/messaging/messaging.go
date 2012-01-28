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

type Message struct {
	Payload []byte
}

func NewMessage(p []byte) (m *Message) {
	return &Message{Payload: p}
}

func NewStringMessage(body string) (m *Message) {
	return NewMessage([]byte(body))
}

type Pull struct {
	Chan   chan *Message
	socket zmq.Socket
}

func NewPull() (p *Pull, err error) {
	var socket zmq.Socket

	socket, err = ctx.NewSocket(zmq.PULL)
	if err != nil {
		return
	}

	p = &Pull{
		Chan:   make(chan *Message),
		socket: socket,
	}

	go func() {
		for {
			body, err := p.socket.Recv(0)
			if err != nil {
				panic(err)
			}

			p.Chan <- &Message{Payload: body}
		}
	}()

	return
}

func (p *Pull) Connect(addr string) (err error) {
	return p.socket.Connect(addr)
}

func (p *Pull) Close() (err error) {
	return p.socket.Close()
}

type Push struct {
	Chan   chan *Message
	socket zmq.Socket
}

func NewPush() (p *Push, err error) {
	var socket zmq.Socket

	socket, err = ctx.NewSocket(zmq.PUSH)
	if err != nil {
		return
	}

	p = &Push{
		Chan:   make(chan *Message),
		socket: socket,
	}

	go func() {
		for {
			msg := <-p.Chan

			// TODO review message buffering
			if err := p.socket.Send(msg.Payload, 0); err != nil {
				panic(err)
			}
		}
	}()

	return
}

func (p *Push) Bind(addr string) (err error) {
	return p.socket.Bind(addr)
}

func (p *Push) Close() (err error) {
	return p.socket.Close()
}
