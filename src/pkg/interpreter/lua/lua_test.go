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
	f, e := l.Eval("local params={...}; return params[1]")
	if e != nil {
		t.Fatalf("Could not compile code: %s", e)
	}
	s := f([]byte("data"), []byte{})
	if string(s) != "data" {
		t.Fatalf("Returned data was not \"data\" but \"%s\"", string(s))
	}
}

