package modifier

import "regexp"

// RequestSleepRule is a rule to delay requests.
type RequestSleepRule struct {
	PathMatcher *regexp.Regexp
	SleepSec    uint
}

// Match tells whether the path matches this rule.
func (rule *RequestSleepRule) Match(path string) bool {
	return rule.PathMatcher.MatchString(path)
}
