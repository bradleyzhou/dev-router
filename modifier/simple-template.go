package modifier

import "strings"

type simpleTemplateKeyword struct {
	// The specific string that serves as a template. For example, "${HOST}" (in a template string "http://${HOST}")
	Key string
	// The value that take the place of Key in the final output. For example, "localhost" (replaces "${HOST}" in "http://${HOST}" to get "http://localhost")
	Value string
}

func writeTemplate(keys []simpleTemplateKeyword, raw string) string {
	v := raw
	for _, key := range keys {
		v = strings.ReplaceAll(v, key.Key, key.Value)
	}
	return v
}
