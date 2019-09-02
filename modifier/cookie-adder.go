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
func (rule *AddRequestCookieRule) AddCookie(host HostDomain, req *http.Request) {
	rule.CookieAdder.Add(host, req)
}

// CookieAdder is the cookie to be added. CookieAdder.Value can contain "${DOMAIN}" for current domain name.
type CookieAdder struct {
	Name  string
	Value string
}

// Add modifies the request by adding a new cookie in the header. Supports dynamic "${DOMAIN}" in cookie value.
func (adder *CookieAdder) Add(host HostDomain, req *http.Request) {
	req.AddCookie(&http.Cookie{Name: adder.Name, Value: simpleHostDomainTemplate(host, adder.Value)})
}
