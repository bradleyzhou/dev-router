package modifier

import "strings"

// HostDomain is a compiled set of domain names of a host.
type HostDomain struct {
	// The host name, e.g. "a.b.example.com:8080"
	Host string
	// The 2nd-level domain name, e.g. "example.com"
	Domain2 string
	// The 3rd-level domain name, e.g. "b.example.com"
	Domain3 string
}

// CompileHostDomain parses the host name and extracts domain names.
// For example, a.b.example.com:8080" --> ("example.com",  "b.example.com")
func CompileHostDomain(host string) HostDomain {
	d2, d3 := extractDomainName(host)
	return HostDomain{
		Host:    host,
		Domain2: d2,
		Domain3: d3,
	}
}

// extractDomainName turns a host name into 2nd-level and 3rd-level domain names.
// For example, "www.a.example.com" --> domain2: "example.com", domain3: "a.example.com"
func extractDomainName(host string) (domain2 string, domain3 string) {
	noPort := strings.Split(host, ":")[0]
	domains := strings.Split(noPort, ".")
	nSubs := len(domains)
	var dn1, dn2, dn3 string
	dn1 = domains[nSubs-1]
	dn2 = dn1
	dn3 = dn2
	if nSubs > 1 {
		dn2 = domains[nSubs-2] + "." + dn1
		dn3 = dn2
	}
	if nSubs > 2 {
		dn3 = domains[nSubs-3] + "." + dn3
	}
	return dn2, dn3
}
