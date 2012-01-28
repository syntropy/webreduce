package messaging

import (
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	addr := "tcp://127.0.0.1:12121"
	testPayload := "teh lulz"

	pub, err := NewPub()
	if err != nil {
		t.Error(err)
	}
	sub, err := NewSub()
	if err != nil {
		t.Error(err)
	}

	err = sub.Connect(addr)
	if err != nil {
		t.Error(err)
	}
	err = pub.Bind(addr)
	if err != nil {
		t.Error(err)
	}

	// FIXME find a better way to detect successfull tcp handshake
	time.Sleep(1 * time.Millisecond)
	pub.Chan <- &Message{Payload: []byte(testPayload)}
	msg := <-sub.Chan
	if string(msg.Payload) != testPayload {
		t.Errorf("expected %s got %s", testPayload, string(msg.Payload))
	}
}

func TestPushPull(t *testing.T) {
	addr := "tcp://127.0.0.1:11111"
	testPayload := "hello"

	pull, err := NewPull()
	if err != nil {
		t.Error(err)
	}
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
