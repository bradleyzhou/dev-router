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
func (rule *PatchHeaderRule) Patch(domain2 string, domain3 string, header http.Header) {
	replacer := writeTemplate([]simpleTemplateKeyword{
		{Key: "${DOMAIN}", Value: domain2},
		{Key: "${DOMAIN_2}", Value: domain2},
		{Key: "${DOMAIN_3}", Value: domain3},
	}, rule.Replacer)

	headers := header[rule.Name]
	for i := range headers {
		headers[i] = rule.Matcher.ReplaceAllString(headers[i], replacer)
	}
}
