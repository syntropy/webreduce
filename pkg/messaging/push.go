package messaging

import (
	zmq "github.com/alecthomas/gozmq"
)

type Push struct {
	Chan   chan *Message
	Err    chan error
	socket zmq.Socket
}

func NewPush(ctx zmq.Context) (p *Push, err error) {
	var socket zmq.Socket

	socket, err = ctx.NewSocket(zmq.PUSH)
	if err != nil {
		return
	}

	p = &Push{
		Chan:   make(chan *Message),
		Err:    make(chan error),
		socket: socket,
	}

	go func() {
		for {
			msg := <-p.Chan

			// TODO review message buffering
			if err := p.socket.Send(msg.Payload, 0); err != nil {
				p.Err <- err
				break
			}
		}
	}()

	return
}

func (p *Push) Bind(addr string) (err error) {
	return p.socket.Bind(addr)
}

func (p *Push) Connect(addr string) (err error) {
	return p.socket.Connect(addr)
}

func (p *Push) Close() (err error) {
	return p.socket.Close()
}
