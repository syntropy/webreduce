package messaging

import (
	"testing"
)

func TestPull(t *testing.T) {
	addr := "tcp://127.0.0.1:11111"
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

	dev.BeginPull("test")

	push.Chan <- &Message{Payload: []byte(testPayload)}
	msg := <-dev.In
	if string(msg.Payload) != testPayload {
		t.Errorf("expected %s got %s", testPayload, string(msg.Payload))
	}
}

func reportDeviceError(t *testing.T, dev *Device) {
	for err := range dev.Err {
		t.Error(err)
	}
}
