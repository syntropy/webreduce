package router

import (
	"net/http"
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

func TestDynamicRuleSpec(t *testing.T) {
	spec := "/test/<foo>/<bar>"
	names := []string{"foo", "bar"}
	values := map[string]string{"foo": "baz", "bar": "qux"}
	regex := ("/test/" + NameGroup + "/" + NameGroup)
	rule := NewRule(spec)

	if vl, nl := len(rule.variables), len(names); vl != nl {
		t.Logf("Got %v variable names, expected %v.", vl, nl)
		t.FailNow()
	}

	for idx, name := range names {
		if v := rule.variables[idx]; v != name {
			t.Errorf("Got variable name '%v' on position %v, expected %v.", v, idx, name)
		}
	}

	if rule.regex.String() != regex {
		t.Errorf("Got Rule.regex %v, expected %v", rule.regex, regex)
	}

	variables, match := rule.Match("/test/"+values["foo"]+"/"+values["bar"], "GET")
	if !match {
		t.Logf("Got match %v, expected %v", match, true)
		t.FailNow()
	}

	for _, name := range names {
		value, found := variables[name]
		if found {
			if v := values[name]; value != v {
				t.Errorf("Got %v for name '%v', expected %v", value, name, v)
			}
		} else {
			t.Errorf("Expected variable '%v' to be in variables.", name)
		}
	}
}

func Handler(w http.ResponseWriter, req *http.Request) {}

func TestStaticRouter(t *testing.T) {
	router := NewRouter()
	patterns := []string{"/foo", "/bar"}

	for _, pattern := range patterns {
		if _, m := router.Match(pattern, "GET"); m {
			t.Errorf("Pattern '%v' shouldn't match", pattern)
		}
	}

	for _, pattern := range patterns {
		router.AddRoute(pattern, Handler)
	}

	for _, pattern := range patterns {
		if _, m := router.Match(pattern, "GET"); !m {
			t.Errorf("Pattern '%v' should match", pattern)
		}
	}
}

func TestDynamicRouter(t *testing.T) {
	pattern := "/<foo>/<bar>"
	path := "/baz/qux"
	router := NewRouter()
	router.AddRoute(pattern, Handler)
	vs, m := router.Match(path, "GET")

	if !m {
		t.Logf("Expected '' to match.", path)
		t.FailNow()
	}

	for k, v := range map[string]string{"foo": "baz", "bar": "qux"} {
		val, found := vs[k]
		if !found {
			t.Errorf("Expected variable '%v' to be in variables.", v)
		} else if val != v {
			t.Errorf("Got '%v' for %v, expected '%v'", val, k, v)
		}

	}
}
