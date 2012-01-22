package router

import (
	"net/http"
	"regexp"
)

const (
	StartToken  = "<"
	EndToken    = ">"
	NameToken   = "[a-z]+"
	NameGroup   = "(" + NameToken + ")"
	NamePattern = StartToken + NameGroup + EndToken
)

// the regular expression to extract variable names from rule specs
var ruleRegex *regexp.Regexp = regexp.MustCompile(NamePattern)

// parses a rule spec in the form /foo/<bar> where 'bar' is the name of a
// variable. It returns a regex.Regexp for extraction of variables and
// the list of names for those variables
func parseRuleSpec(spec string) (regex *regexp.Regexp, vs []string) {
	matches := ruleRegex.FindAllStringSubmatch(spec, -1)

	if matches == nil {
		regex = regexp.MustCompile(regexp.QuoteMeta(spec))
		return
	}

	for _, match := range matches {
		vs = append(vs, match[1])
		spec = regexp.MustCompile(match[0]).ReplaceAllString(spec, "([a-z]+)")
	}

	regex = regexp.MustCompile(spec)

	return
}

// Rule is the representation of a URL rule
type Rule struct {
	spec      string
	methods   []string
	regex     *regexp.Regexp
	variables []string
}

// Create a new Rule. The rule spec ...
// If 'GET' is provided as methods, 'HEAD' is automaticaly added
// to the new rule.
func NewRule(spec string, methods ...string) *Rule {
	regex, variables := parseRuleSpec(spec)

	verbs := map[string]bool{}
	ms := []string{}

	if len(methods) != 0 {
		for _, method := range methods {
			switch method {
			case "GET":
				verbs["HEAD"] = true
				verbs["GET"] = true
			default:
				verbs[method] = true
			}
		}

		for method, _ := range verbs {
			ms = append(ms, method)
		}

	} else {
		ms = []string{"HEAD", "GET"}
	}

	return &Rule{spec, ms, regex, variables}
}

// Match given path and return named variables.
func (r *Rule) Match(pattern string, method string) (variables map[string]string, match bool) {
	for _, m := range r.methods {
		if m == method {
			match = true
			break
		}
	}

	if !match {
		return
	}

	match = r.regex.MatchString(pattern)

	if !match {
		return
	}

	values := r.regex.FindStringSubmatch(pattern)
	if values == nil {
		return
	}

	variables = map[string]string{}
	for i, v := range values[1:] {
		variables[r.variables[i]] = v
	}

	return
}

// A Route associates a Rule and a handler function
type Route struct {
	rule Rule
	handler func(http.ResponseWriter, *http.Request)
}

// Match given pattern. Documented in Rule.Match
func (r *Route) Match(pattern string, method string) (variables map[string]string, match bool) {
	return r.rule.Match(pattern, method)
}

// A Router dispatches HTTP requests to handlers.
type Router struct {
	routes []Route
}

// Create a new Router
func NewRouter() Router {
	r := Router{}

	return r
}

// Add a route to this router.
func (r *Router) AddRoute(spec string, handler func(http.ResponseWriter, *http.Request), methods ...string) {
	rule := *NewRule(spec, methods...)
	route := Route{rule, handler}

	r.routes = append(r.routes, route)
}

// Match given pattern. Documented in Rule.Match
func (r *Router) Match(pattern string, method string) (match bool) {
	for _, route := range r.routes {
		if _, matched := route.Match(pattern, method); matched {
			return matched
		}
	}
	return
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}
