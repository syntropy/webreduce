package messaging

import (
	zmq "github.com/alecthomas/gozmq"
)

type Pub struct {
	Chan   chan *Message
	socket zmq.Socket
}

func NewPub() (p *Pub, err error) {
	var socket zmq.Socket

	socket, err = ctx.NewSocket(zmq.PUB)
	if err != nil {
		return
	}

	p = &Pub{
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

func (p *Pub) Bind(addr string) (err error) {
	return p.socket.Bind(addr)
}

func (p *Pub) Close() (err error) {
	return p.socket.Close()
}
