package router

import (
	// "fmt"
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
func (r *Rule) MatchPath(path string) (variables map[string]string, match bool) {
	match = r.regex.MatchString(path)

	if !match {
		return
	}

	values := r.regex.FindStringSubmatch(path)
	if values == nil {
		return
	}

	variables = map[string]string{}
	for i, v := range values[1:] {
		variables[r.variables[i]] = v
	}

	return
}
