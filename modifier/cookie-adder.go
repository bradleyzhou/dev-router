package modifier

import (
	"net/http"
	"regexp"
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
	v := writeTemplate([]simpleTemplateKeyword{
		{Key: "${DOMAIN}", Value: domain2},
		{Key: "${DOMAIN_2}", Value: domain2},
		{Key: "${DOMAIN_3}", Value: domain3},
	}, adder.Value)
	req.AddCookie(&http.Cookie{Name: adder.Name, Value: v})
}
