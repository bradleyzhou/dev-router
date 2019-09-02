package modifier

import "strings"

type simpleTemplateKeyword struct {
	// The specific string that serves as a template. For example, "${HOST}" (in a template string "http://${HOST}")
	Key string
	// The value that take the place of Key in the final output. For example, "localhost" (replaces "${HOST}" in "http://${HOST}" to get "http://localhost")
	Value string
}

// writeTemplate is a rudimentary template engine.
func writeTemplate(keys []simpleTemplateKeyword, raw string) string {
	v := raw
	for _, key := range keys {
		v = strings.ReplaceAll(v, key.Key, key.Value)
	}
	return v
}

// simpleHostDomainTemplate replaces ${HOST} and ${DOMAIN} with a provided host and domain.
// Useful for dynamically replace the host name and domain name.
func simpleHostDomainTemplate(host HostDomain, raw string) string {
	return writeTemplate([]simpleTemplateKeyword{
		{Key: "${HOST}", Value: host.Host},
		{Key: "${DOMAIN}", Value: host.Domain2},
		{Key: "${DOMAIN_2}", Value: host.Domain2},
		{Key: "${DOMAIN_3}", Value: host.Domain3},
	}, raw)
}
