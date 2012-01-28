package messaging

import (
	"testing"
)

func TestPushPull(t *testing.T) {
	addr := "tcp://127.0.0.1:11111"
	testPayload := "hello"
	pull, err := NewPull()
	push, err := NewPush()
	if err != nil {
		t.Error(err)
	}

	err = push.Bind(addr)
	if err != nil {
		t.Error(err)
	}
	err = pull.Connect(addr)
	if err != nil {
		t.Error(err)
	}

	push.Chan <- &Message{Payload: []byte(testPayload)}
	msg := <-pull.Chan
	if string(msg.Payload) != testPayload {
		t.Errorf("expected %s got %s", testPayload, string(msg.Payload))
	}
}
