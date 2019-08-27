package modifier

import (
	"net/http"
)

// AddHeaderRule TODO
type AddHeaderRule struct {
	Name  string
	Value string
}

// Add TODO
func (rule *AddHeaderRule) Add(domain2 string, domain3 string, header http.Header) {
	v := writeTemplate([]simpleTemplateKeyword{
		{Key: "${DOMAIN}", Value: domain2},
		{Key: "${DOMAIN_2}", Value: domain2},
		{Key: "${DOMAIN_3}", Value: domain3},
	}, rule.Value)
	header.Add(rule.Name, v)
}
