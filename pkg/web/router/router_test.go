package router

import (
	"net/http"
	"testing"
)

func TestStaticRulePattern(t *testing.T) {
	pattern := "/test"
	rule := NewRule(pattern)

	if rule.pattern != pattern {
		t.Errorf("Got Rule.pattern %v, expected %v", rule.pattern, pattern)
	}
}

func TestStaticRuleHeadInjection(t *testing.T) {
	pattern := "/test"
	rule := NewRule(pattern)

	if l := len(rule.methods); l != 2 {
		t.Errorf("Got Rule.methods of length %v, expected 2", l)
	}

	methods := map[string]bool{"GET": true, "HEAD": true}

	for _, method := range rule.methods {
		if _, found := methods[method]; !found {
			t.Errorf("Got unexpected method '%v'", method)
		}
	}
}

func TestStaticRuleWithMethods(t *testing.T) {
	pattern := "/test"
	rule := NewRule(pattern, "GET", "PUT", "POST", "DELETE")
	methods := map[string]bool{
		"HEAD":   true,
		"GET":    true,
		"PUT":    true,
		"POST":   true,
		"DELETE": true,
	}

	for _, method := range rule.methods {
		if _, found := methods[method]; !found {
			t.Errorf("Got unexpected method '%v'", method)
		}
	}
}

func TestDynamicRulePattern(t *testing.T) {
	pattern := "/test/<foo>/<bar>"
	names := []string{"foo", "bar"}
	values := map[string]string{"foo": "baz", "bar": "qux"}
	regex := ("^/test/" + NameGroup + "/" + NameGroup + "$")
	rule := NewRule(pattern)

	if vl, nl := len(rule.args), len(names); vl != nl {
		t.Logf("Got %v arg names, expected %v.", vl, nl)
		t.FailNow()
	}

	for idx, name := range names {
		if v := rule.args[idx]; v != name {
			t.Errorf("Got arg name '%v' on position %v, expected %v.", v, idx, name)
		}
	}

	if rule.regex.String() != regex {
		t.Errorf("Got Rule.regex %v, expected %v", rule.regex, regex)
	}

	args, match := rule.Match("/test/"+values["foo"]+"/"+values["bar"], "GET")
	if !match {
		t.Logf("Got match %v, expected %v", match, true)
		t.FailNow()
	}

	for _, name := range names {
		value, found := args[name]
		if found {
			if v := values[name]; value != v {
				t.Errorf("Got %v for name '%v', expected %v", value, name, v)
			}
		} else {
			t.Errorf("Expected arg '%v' to be in args.", name)
		}
	}
}

func Handler(ctx map[string]string, w http.ResponseWriter, req *http.Request) {}

func TestStaticRouter(t *testing.T) {
	prefix := "/qux"
	router := NewRouter(prefix)
	patterns := []string{"/foo", "/bar"}

	for _, pattern := range patterns {
		if _, _, m := router.Match(prefix+pattern, "GET"); m {
			t.Errorf("Pattern '%v' shouldn't match", pattern)
		}
	}

	for _, pattern := range patterns {
		router.AddRoute(pattern, Handler)
	}

	for _, pattern := range patterns {
		if _, _, m := router.Match(prefix+pattern, "GET"); !m {
			t.Errorf("Pattern '%v' should match", pattern)
		}
	}
}

func TestDynamicRouter(t *testing.T) {
	pattern := "/<foo>/<bar>"
	path := "/baz/qux"
	router := NewRouter("")
	router.AddRoute(pattern, Handler)
	_, vs, m := router.Match(path, "GET")

	if !m {
		t.Logf("Expected '%v' to match.", path)
		t.FailNow()
	}

	for k, v := range map[string]string{"foo": "baz", "bar": "qux"} {
		val, found := vs[k]
		if !found {
			t.Errorf("Expected arg '%v' to be in args.", v)
		} else if val != v {
			t.Errorf("Got '%v' for %v, expected '%v'", val, k, v)
		}
	}
}
