package modifier

import (
	"net/http"
)

// AddHeaderRule represents a header entry. Value supports dynamic "${DOMAIN}".
type AddHeaderRule struct {
	Name  string
	Value string
}

// Add modifies the header by adding a new header entry. Supports dynamic "${DOMAIN}".
func (rule *AddHeaderRule) Add(host HostDomain, header http.Header) {
	header.Add(rule.Name, simpleHostDomainTemplate(host, rule.Value))
}
