package messaging

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"strconv"
)

type Device struct {
	In      chan *Message
	Out     chan *Message
	Err     chan error
	sockets map[string]map[string]Socket
	ctx     zmq.Context
}

func NewDevice() (d *Device, err error) {
	ctx, err := zmq.NewContext()
	if err != nil {
		return
	}

	d = &Device{
		In:      make(chan *Message),
		Out:     make(chan *Message),
		Err:     make(chan error),
		sockets: map[string]map[string]Socket{},
		ctx:     ctx,
	}

	return
}

func (d *Device) StartPub(endpoint string) {
	id := "pub/" + endpoint

	d.sockets[id] = map[string]Socket{}

	// XXX mocks a potential return value from the coordinator
	addrs := []string{"ipc:///tmp/pub"}

	for i := range addrs {
		d.addPub(id, addrs[i])
	}
}

func (d *Device) StopPub(endpoint string) {
	id := "pub/" + endpoint

	for key, _ := range d.sockets[id] {
		d.deleteSocket(id, key)
	}

	delete(d.sockets, endpoint)
}

func (d *Device) StartPull(endpoint string) {
	id := "pull/" + endpoint

	d.sockets[id] = map[string]Socket{}

	// XXX mocks a potential return value from the coordinator
	addrs := []string{"ipc:///tmp/push"}

	for i := range addrs {
		d.addPull(id, addrs[i])
	}
}

func (d *Device) StopPull(endpoint string) {
	id := "pull/" + endpoint

	for key, _ := range d.sockets[id] {
		d.deleteSocket(id, key)
	}

	delete(d.sockets, endpoint)
}

func (d *Device) StartSub(endpoint string) {
	id := "sub/" + endpoint

	d.sockets[id] = map[string]Socket{}

	// XXX mocks a potential return value from the coordinator
	addrs := []string{"ipc:///tmp/pub"}

	for i := range addrs {
		d.addSub(id, addrs[i])
	}
}

func (d *Device) StopSub(endpoint string) {
	id := "sub/" + endpoint

	for key, _ := range d.sockets[id] {
		d.deleteSocket(id, key)
	}

	delete(d.sockets, endpoint)
}

func (d *Device) String() string {
	return fmt.Sprintf("%#v", d)
}

func (d *Device) addPub(id string, addr string) {
	socks := d.sockets[id]
	k := strconv.Itoa(len(socks) + 1)
	p, err := NewPub(d.ctx)
	d.emitError(err)

	go func() {
		for {
			err := <-p.Err
			d.emitError(err)

			d.deleteSocket(id, k)

			break
		}
	}()

	err = p.Bind(addr)
	d.emitError(err)

	go func() {
		for msg := range d.Out {
			p.Chan <- msg
		}
	}()

	socks[k] = p
}

func (d *Device) addPull(id string, addr string) {
	socks := d.sockets[id]
	k := strconv.Itoa(len(socks) + 1)
	p, err := NewPull(d.ctx)
	d.emitError(err)

	go func() {
		for {
			err := <-p.Err
			d.emitError(err)

			d.deleteSocket(id, k)

			break
		}
	}()

	err = p.Connect(addr)
	d.emitError(err)

	go func() {
		for msg := range p.Chan {
			d.In <- msg
		}
	}()

	socks[k] = p
}

func (d *Device) addSub(id string, addr string) {
	socks := d.sockets[id]
	k := strconv.Itoa(len(socks) + 1)
	sock, err := NewSub(d.ctx)
	d.emitError(err)

	go func() {
		for {
			err := <-sock.Err
			d.emitError(err)

			d.deleteSocket(id, k)

			break
		}
	}()

	err = sock.Connect(addr)
	d.emitError(err)

	go func() {
		for msg := range sock.Chan {
			d.In <- msg
		}
	}()

	socks[k] = sock
}

func (d *Device) deleteSocket(id string, key string) {
	socks := d.sockets[id]

	err := socks[key].Close()
	d.emitError(err)

	delete(socks, key)
}

func (d *Device) emitError(err error) {
	if err != nil {
		d.Err <- err
	}
}
