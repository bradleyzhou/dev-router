package modifier

import (
	"net/http"
	"regexp"
	"strings"
)

// PatchHeaderRule TODO
type PatchHeaderRule struct {
	Name     string
	Matcher  *regexp.Regexp
	Replacer string
}

// Patch TODO
func (rule *PatchHeaderRule) Patch(domain2 string, domain3 string, header http.Header) {
	replacer := strings.ReplaceAll(rule.Replacer, "${DOMAIN}", domain2)
	replacer = strings.ReplaceAll(replacer, "${DOMAIN_2}", domain2)
	replacer = strings.ReplaceAll(replacer, "${DOMAIN_3}", domain3)
	headers := header[rule.Name]
	for i := range headers {
		headers[i] = rule.Matcher.ReplaceAllString(headers[i], replacer)
	}
}
