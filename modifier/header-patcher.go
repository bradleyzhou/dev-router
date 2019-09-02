package modifier

import (
	"net/http"
	"regexp"
)

// PatchHeaderRule represents a rule for modify an existing header.
type PatchHeaderRule struct {
	Name     string
	Matcher  *regexp.Regexp
	Replacer string
}

// Patch modifies the header entries that has the same name in the rule, and
// replace the value according to the rule.
func (rule *PatchHeaderRule) Patch(host HostDomain, header http.Header) {
	headers := header[rule.Name]
	for i := range headers {
		headers[i] = rule.Matcher.ReplaceAllString(headers[i], simpleHostDomainTemplate(host, rule.Replacer))
	}
}
