package messaging

import (
	zmq "github.com/alecthomas/gozmq"
)

type Sub struct {
	Chan   chan *Message
	Err    chan error
	socket zmq.Socket
	prefix string
}

func NewSub(ctx zmq.Context, qName string) (s *Sub, err error) {
	var socket zmq.Socket

	socket, err = ctx.NewSocket(zmq.SUB)
	if err != nil {
		return
	}

	s = &Sub{
		Chan:   make(chan *Message),
		Err:    make(chan error),
		socket: socket,
		prefix: qName,
	}

	err = s.socket.SetSockOptString(zmq.SUBSCRIBE, qName)
	if err != nil {
		return
	}

	go func() {
		for {
			body, err := s.socket.Recv(0)
			if err != nil {
				s.Err <- err
				break
			}

			s.Chan <- &Message{Payload: body[len(s.prefix)+1:]}
		}
	}()

	return
}

func (s *Sub) Bind(addr string) (err error) {
	return s.socket.Bind(addr)
}

func (s *Sub) Connect(addr string) (err error) {
	return s.socket.Connect(addr)
}

func (s *Sub) Close() (err error) {
	return s.socket.Close()
}
