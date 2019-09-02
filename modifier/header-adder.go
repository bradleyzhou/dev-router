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
func (rule *AddHeaderRule) Add(domain2 string, domain3 string, header http.Header) {
	v := writeTemplate([]simpleTemplateKeyword{
		{Key: "${DOMAIN}", Value: domain2},
		{Key: "${DOMAIN_2}", Value: domain2},
		{Key: "${DOMAIN_3}", Value: domain3},
	}, rule.Value)
	header.Add(rule.Name, v)
}
