package router

import (
	"net/http"
	"regexp"
	"sort"
)

const (
	StartToken  = "<"
	EndToken    = ">"
	NameToken   = "[a-z]+"
	NameGroup   = "(" + NameToken + ")"
	NamePattern = StartToken + NameGroup + EndToken
)

// the regular expression to extract arg names from rule patterns
var ruleRegex *regexp.Regexp = regexp.MustCompile(NamePattern)

// parses a rule pattern in the form /foo/<bar> where 'bar' is the name of a
// arg. It returns a regex.Regexp for extraction of args and
// the list of names for those args
func parseRulePattern(pattern string) (regex *regexp.Regexp, vs []string) {
	matches := ruleRegex.FindAllStringSubmatch(pattern, -1)

	if matches == nil {
		regex = regexp.MustCompile("^" + regexp.QuoteMeta(pattern) + "$")
		return
	}

	for _, match := range matches {
		vs = append(vs, match[1])
		pattern = regexp.MustCompile(match[0]).ReplaceAllString(pattern, "([a-z]+)")
	}

	regex = regexp.MustCompile("^" + pattern + "$")

	return
}

// Rule is the representation of a URL rule
type Rule struct {
	pattern string
	methods []string
	regex   *regexp.Regexp
	args    []string
}

// Create a new Rule. The rule pattern ...
// If 'GET' is provided as methods, 'HEAD' is automaticaly added
// to the new rule.
func NewRule(pattern string, methods ...string) *Rule {
	regex, args := parseRulePattern(pattern)

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

	return &Rule{pattern, ms, regex, args}
}

// Match given path and return named args.
func (r *Rule) Match(pattern string, method string) (args map[string]string, match bool) {
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

	args = map[string]string{}
	for i, v := range values[1:] {
		args[r.args[i]] = v
	}

	return
}

// A Route associates a Rule and a handler function
type Route struct {
	rule    Rule
	handler func(map[string]string, http.ResponseWriter, *http.Request)
}

// Match given pattern. Documented in Rule.Match
func (r *Route) Match(pattern string, method string) (args map[string]string, match bool) {
	return r.rule.Match(pattern, method)
}

// A RoutesList is a sortable list of Routes.
type RouteList []Route

// Len is required by sort.Interface.
func (rl RouteList) Len() int {
	return len(rl)
}

// Less is required by sort.Interface.
func (rl RouteList) Less(i, j int) bool {
	return rl[i].rule.pattern < rl[j].rule.pattern
}

// Swap is required by sort.Interface.
func (rl RouteList) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}

// A Router dispatches HTTP requests to handlers.
type Router struct {
	prefix string
	routes RouteList
	sorted bool
}

// Create a new Router
func NewRouter(prefix string) Router {
	r := Router{prefix: prefix}

	return r
}

// Add a route to this router.
func (r *Router) AddRoute(pattern string, handler func(map[string]string, http.ResponseWriter, *http.Request), methods ...string) {
	rule := *NewRule(r.prefix+pattern, methods...)
	route := Route{rule, handler}

	r.routes = append(r.routes, route)
}

// Match given pattern. Documented in Rule.Match
func (r *Router) Match(pattern string, method string) (handler func(map[string]string, http.ResponseWriter, *http.Request), args map[string]string, match bool) {
	if !r.sorted {
		sort.Sort(r.routes)
		r.sorted = true
	}

	for _, route := range r.routes {
		if args, matched := route.Match(pattern, method); matched {
			return route.handler, args, matched
		}
	}
	return
}

// Implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, args, matched := r.Match(req.URL.Path, req.Method)
	if !matched {
		http.NotFound(w, req)
		return
	}

	handler(args, w, req)
}
