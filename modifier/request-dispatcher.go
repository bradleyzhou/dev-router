package modifier

import (
	"regexp"
)

// DirectedLocation TODO
type DirectedLocation struct {
	Server string
	Scheme string
	Host   string
	Path   string
}

// RequestDispatchRule TODO
type RequestDispatchRule struct {
	PathMatcher  *regexp.Regexp
	PathReplacer string
	DstScheme    string
	DstHost      string
	DstServer    string
}

// Match TODO
func (rule *RequestDispatchRule) Match(path string) bool {
	return rule.PathMatcher.MatchString(path)
}

// Direct TODO
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
