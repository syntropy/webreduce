package messaging

import (
	zmq "github.com/alecthomas/gozmq"
)

type Pub struct {
	Chan   chan *Message
	Err    chan error
	socket zmq.Socket
	prefix string
}

func NewPub(ctx zmq.Context, qName string) (p *Pub, err error) {
	var socket zmq.Socket

	socket, err = ctx.NewSocket(zmq.PUB)
	if err != nil {
		return
	}

	p = &Pub{
		Chan:   make(chan *Message),
		Err:    make(chan error),
		socket: socket,
		prefix: qName,
	}

	go func() {
		for {
			msg := <-p.Chan
			payload := append([]byte(p.prefix+" "), msg.Payload...)

			// TODO review message buffering
			if err := p.socket.Send(payload, 1); err != nil {
				p.Err <- err
				break
			}
		}
	}()

	return
}

func (p *Pub) Bind(addr string) (err error) {
	return p.socket.Bind(addr)
}

func (p *Pub) Connect(addr string) (err error) {
	return p.socket.Connect(addr)
}

func (p *Pub) Close() (err error) {
	return p.socket.Close()
}
