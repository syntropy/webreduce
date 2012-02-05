package messaging

import (
	"testing"
	"time"
)

func TestPub(t *testing.T) {
	addr := "tcp://127.0.0.1:11111"
	testPayload := "pub device test"
	dev := NewDevice()
	go reportDeviceError(t, dev)

	sub, err := NewSub()
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

	dev.StopPub("pub-test")
}

func TestPull(t *testing.T) {
	addr := "tcp://127.0.0.1:11112"
	testPayload := "hello"
	dev := NewDevice()
	go reportDeviceError(t, dev)

	push, err := NewPush()
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
}

func reportDeviceError(t *testing.T, dev *Device) {
	for err := range dev.Err {
		t.Error(err)
	}
}
