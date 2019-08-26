package modifier

import (
	"net/http"
	"regexp"
	"strings"
)

// AddRequestCookieRule TOOD
type AddRequestCookieRule struct {
	PathMatcher *regexp.Regexp
	CookieAdder CookieAdder
}

// Match TODO
func (rule *AddRequestCookieRule) Match(path string) bool {
	return rule.PathMatcher.MatchString(path)
}

// AddCookie TODO
func (rule *AddRequestCookieRule) AddCookie(domain2 string, domain3 string, req *http.Request) {
	rule.CookieAdder.Add(domain2, domain3, req)
}

// CookieAdder is the cookie to be added. CookieAdder.Value can contain "${DOMAIN}" for current domain name
type CookieAdder struct {
	Name  string
	Value string
}

// Add TODO
func (adder *CookieAdder) Add(domain2 string, domain3 string, req *http.Request) {
	v := strings.ReplaceAll(adder.Value, "${DOMAIN}", domain2)
	v = strings.ReplaceAll(v, "${DOMAIN_2}", domain2)
	v = strings.ReplaceAll(v, "${DOMAIN_3}", domain3)
	req.AddCookie(&http.Cookie{Name: adder.Name, Value: v})
}
