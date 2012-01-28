package lua

import (
	"testing"
)

func TestEmitting(t *testing.T) {
	l := New()
	f, e := l.Eval("a = 1 + 3; emit(a); return a+1")
	if e != nil {
		t.Fatalf("Could not compile code: %s", e)
	}
	l.RegisterEmitCallback(func(data []byte) {
		if string(data) != "4" {
			t.Fatalf("Emitted data was not 4 but \"%s\"", string(data))
		}
	})
	s := f([]byte{}, []byte{})
	if string(s) != "5" {
		t.Fatalf("Returned data was not 9 but \"%s\"", string(s))
	}
}

func TestDataPassing(t *testing.T) {
	l := New()
	f, e := l.Eval("local params={...}; return params[1] .. params[2];")
	if e != nil {
		t.Fatalf("Could not compile code: %s", e)
	}
	data, state := "data", "state"
	s := f([]byte(data), []byte(state))
	if string(s) != data+state {
		t.Fatalf("Returned data was not \"%s\" but \"%s\"", data+state, string(s))
	}
}

func TestInvalidCode(t *testing.T) {
	l := New()
	_, e := l.Eval("This is not code")
	if e == nil {
		t.Fatalf("Could compile code")
	}
}
