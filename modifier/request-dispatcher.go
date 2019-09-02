package modifier

import (
	"regexp"
)

// DirectedLocation represents the target location where the router relays to.
type DirectedLocation struct {
	Server string
	Scheme string
	Host   string
	Path   string
}

// RequestDispatchRule is a rule to dispatch specific requests.
type RequestDispatchRule struct {
	PathMatcher  *regexp.Regexp
	PathReplacer string
	DstScheme    string
	DstHost      string
	DstServer    string
}

// Match tells whether the path matches this rule.
func (rule *RequestDispatchRule) Match(path string) bool {
	return rule.PathMatcher.MatchString(path)
}

// Direct turns a path into an appropriate DirectedLocation for later consumption.
// Supports the dynamic "${PATH}" in target path.
func (rule *RequestDispatchRule) Direct(path string) DirectedLocation {
	// rudimentary support for directing to a dynamic path if replacer is "${PATH}"
	var newPath string
	if rule.PathReplacer == "${PATH}" {
		newPath = path
	} else {
		newPath = rule.PathMatcher.ReplaceAllString(path, rule.PathReplacer)
	}

	return DirectedLocation{
		Server: rule.DstServer,
		Scheme: rule.DstScheme,
		Host:   rule.DstHost,
		Path:   newPath,
	}
}
