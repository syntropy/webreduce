package router

import (
	"testing"
)

func TestStaticRuleSpec(t *testing.T) {
	spec := "/test"
	rule := NewRule(spec)

	if rule.spec != spec {
		t.Errorf("Got Rule.spec %v, expected %v", rule.spec, spec)
	}
}

func TestStaticRuleHeadInjection(t *testing.T) {
	spec := "/test"
	rule := NewRule(spec)

	if l := len(rule.methods); l != 2 {
		t.Errorf("Got Rule.methods of length %v, expected 2", l)
	}

	methods := map[string]bool{"GET": true, "HEAD": true}

	for _, method := range rule.methods {
		if _, found := methods[method]; found != true {
			t.Errorf("Got unexpected method '%v'", method)
		}
	}
}

func TestStaticRuleWithMethods(t *testing.T) {
	spec := "/test"
	rule := NewRule(spec, "GET", "PUT", "POST", "DELETE")
	methods := map[string]bool{
		"HEAD":   true,
		"GET":    true,
		"PUT":    true,
		"POST":   true,
		"DELETE": true,
	}

	for _, method := range rule.methods {
		if _, found := methods[method]; found != true {
			t.Errorf("Got unexpected method '%v'", method)
		}
	}
}
