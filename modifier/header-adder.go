package modifier

import (
	"net/http"
	"strings"
)

// AddHeaderRule TODO
type AddHeaderRule struct {
	Name  string
	Value string
}

// Add TODO
func (rule *AddHeaderRule) Add(domain2 string, domain3 string, header http.Header) {
	v := strings.ReplaceAll(rule.Value, "${DOMAIN}", domain2)
	v = strings.ReplaceAll(v, "${DOMAIN_2}", domain2)
	v = strings.ReplaceAll(v, "${DOMAIN_3}", domain3)
	header.Add(rule.Name, v)
}
