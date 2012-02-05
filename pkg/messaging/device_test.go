package messaging

import (
	zmq "github.com/alecthomas/gozmq"
	"testing"
	"time"
)

func TestPub(t *testing.T) {
	addr := "ipc:///tmp/pub-test"
	testPayload := "pub device test"
	dev, err := NewDevice()
	if err != nil {
		t.Error(err)
	}
	go reportDeviceError(t, dev)

	sub, err := NewSub(createContext(), "pub-test")
	if err != nil {
		t.Error(err)
	}
	err = sub.Connect(addr)
	if err != nil {
		t.Error(err)
	}

	dev.StartPub("pub-test")

	// FIXME find a better way to detect successfull tcp handshake
	time.Sleep(200 * time.Millisecond)
	dev.Out <- &Message{Payload: []byte(testPayload)}
	msg := <-sub.Chan
	if string(msg.Payload) != testPayload {
		t.Errorf("expected %s got %s", testPayload, string(msg.Payload))
	}
}

func TestPull(t *testing.T) {
	addr := "ipc:///tmp/pull-test"
	testPayload := "hello"
	dev, err := NewDevice()
	if err != nil {
		t.Error(err)
	}
	go reportDeviceError(t, dev)

	push, err := NewPush(createContext())
	if err != nil {
		t.Error(err)
	}
	err = push.Bind(addr)
	if err != nil {
		t.Error(err)
	}

	dev.StartPull("pull-test")

	push.Chan <- &Message{Payload: []byte(testPayload)}
	msg := <-dev.In
	if string(msg.Payload) != testPayload {
		t.Errorf("expected %s got %s", testPayload, string(msg.Payload))
	}

	dev.StopPull("pull-test")
	push.Close()
}

func TestSub(t *testing.T) {
	addr := "ipc:///tmp/sub-test"
	testPayload := "payloadz"
	dev, err := NewDevice()
	if err != nil {
		t.Error(err)
	}
	go reportDeviceError(t, dev)

	pub, err := NewPub(createContext(), "sub-test")
	if err != nil {
		t.Error(err)
	}
	err = pub.Bind(addr)
	if err != nil {
		t.Error(err)
	}

	dev.StartSub("sub-test")

	// FIXME find a better way to detect successfull tcp handshake
	time.Sleep(200 * time.Millisecond)
	pub.Chan <- &Message{Payload: []byte(testPayload)}
	msg := <-dev.In
	if string(msg.Payload) != testPayload {
		t.Errorf("expected %s got %s", testPayload, string(msg.Payload))
	}
}

func createContext() (ctx zmq.Context) {
	ctx, err := zmq.NewContext()
	if err != nil {
		panic(err)
	}

	return
}

func reportDeviceError(t *testing.T, dev *Device) {
	for err := range dev.Err {
		t.Error(err)
	}
}
