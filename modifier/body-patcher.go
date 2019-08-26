package modifier

import (
	"bytes"
	"regexp"
)

// PatchBodyRule represents a rule to modify a (response) body
type PatchBodyRule struct {
	// Matcher is the regex pattern to be replaced
	Matcher *regexp.Regexp
	// Replacer is the replacement text repl (the same as in regex.ReplaceAll),
	// supports "${HOST}" for current host name
	Replacer []byte
}

// ReplaceAll will search and replace the response body according to the rule.
func (rule *PatchBodyRule) ReplaceAll(host string, body []byte) []byte {
	replacer := bytes.ReplaceAll(rule.Replacer, []byte("${HOST}"), []byte(host))
	return rule.Matcher.ReplaceAllLiteral(body, replacer)
}
