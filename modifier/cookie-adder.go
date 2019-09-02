package modifier

import (
	"net/http"
	"regexp"
)

// AddRequestCookieRule represents a rule for adding a new (request) cookie.
type AddRequestCookieRule struct {
	PathMatcher *regexp.Regexp
	CookieAdder CookieAdder
}

// Match tests whether the path matches this rule.
func (rule *AddRequestCookieRule) Match(path string) bool {
	return rule.PathMatcher.MatchString(path)
}

// AddCookie modifies the request by adding the cookie specified in this rule.
// Supports the dynamic rule "${DOMAIN}", so a 2nd-level and 3rd-level domain name is needed.
func (rule *AddRequestCookieRule) AddCookie(domain2 string, domain3 string, req *http.Request) {
	rule.CookieAdder.Add(domain2, domain3, req)
}

// CookieAdder is the cookie to be added. CookieAdder.Value can contain "${DOMAIN}" for current domain name.
type CookieAdder struct {
	Name  string
	Value string
}

// Add modifies the request by adding a new cookie in the header. Supports dynamic "${DOMAIN}" in cookie value.
func (adder *CookieAdder) Add(domain2 string, domain3 string, req *http.Request) {
	v := writeTemplate([]simpleTemplateKeyword{
		{Key: "${DOMAIN}", Value: domain2},
		{Key: "${DOMAIN_2}", Value: domain2},
		{Key: "${DOMAIN_3}", Value: domain3},
	}, adder.Value)
	req.AddCookie(&http.Cookie{Name: adder.Name, Value: v})
}
