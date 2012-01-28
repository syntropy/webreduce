package messaging

import (
	zmq "github.com/alecthomas/gozmq"
)

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
