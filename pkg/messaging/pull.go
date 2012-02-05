package messaging

import (
	zmq "github.com/alecthomas/gozmq"
)

type Pull struct {
	Chan   chan *Message
	Err    chan error
	socket zmq.Socket
}

func NewPull(ctx zmq.Context) (p *Pull, err error) {
	var socket zmq.Socket

	socket, err = ctx.NewSocket(zmq.PULL)
	if err != nil {
		return
	}

	p = &Pull{
		Chan:   make(chan *Message),
		Err:    make(chan error),
		socket: socket,
	}

	go func() {
		for {
			body, err := p.socket.Recv(0)
			if err != nil {
				p.Err <- err
				break
			}

			p.Chan <- &Message{Payload: body}
		}
	}()

	return
}

func (p *Pull) Bind(addr string) (err error) {
	return p.socket.Bind(addr)
}

func (p *Pull) Connect(addr string) (err error) {
	return p.socket.Connect(addr)
}

func (p *Pull) Close() (err error) {
	return p.socket.Close()
}
