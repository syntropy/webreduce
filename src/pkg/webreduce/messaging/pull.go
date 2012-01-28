package messaging

import (
	zmq "github.com/alecthomas/gozmq"
)

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
