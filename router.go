package main

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/bradleyzhou/dev-router/config"
	"github.com/bradleyzhou/dev-router/configure"
	"github.com/bradleyzhou/dev-router/modifier"
)

var serverConf configure.Config

func patchResponseBody(host string, body []byte, rules []modifier.PatchBodyRule) []byte {
	for _, r := range rules {
		body = r.ReplaceAll(host, body)
	}
	return body
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

func serveReverseProxy(res http.ResponseWriter, req *http.Request) {
	originReqHost := req.Host
	originReqDomain2, originReqDomain3 := extractDomainName(originReqHost)

	director := func(req *http.Request) {
		log.Printf("----->>> Request received.")

		directedScheme := serverConf.DefaultRequestDispatchRule.DstScheme
		directedHost := serverConf.DefaultRequestDispatchRule.DstHost
		directedServer := directedHost

		for _, r := range serverConf.RequestDispatchRules {
			if r.Match(req.URL.Path) {
				dst := r.Direct(req.URL.Path)

				switch dst.Scheme {
				case "":
					// use the scheme in the default rule
				case "${SCHEME}":
					directedScheme = serverConf.ServerScheme
				default:
					directedScheme = dst.Scheme
				}

				switch dst.Host {
				case "":
					// use the host in the default rule
				case "${HOST}":
					directedHost = originReqHost
				default:
					directedHost = dst.Host
				}

				switch dst.Server {
				case "":
					// use the server in the default rule
					directedServer = directedHost
				case "${HOST}":
					directedServer = originReqHost
				default:
					directedServer = dst.Server
				}

				if dst.Path != "" && dst.Path != "/" {
					req.URL.Path = dst.Path
				}
				break
			}
		}

		for _, r := range serverConf.ReqAddCookieRules {
			if r.Match(req.URL.Path) {
				r.AddCookie(originReqDomain2, originReqDomain3, req)
			}
		}

		req.URL.Scheme = directedScheme
		req.URL.Host = directedServer
		req.Host = directedHost

		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		if _, ok := req.Header["Origin"]; ok {
			// mod origin
			req.Header.Set("Origin", directedHost)
		}
		if _, ok := req.Header["Referer"]; ok {
			// mod referer
			req.Header.Set("Referer", directedHost)
		}

		log.Printf("--> Direct using %v to server %v (host %v) path %v", req.URL.Scheme, req.URL.Host, req.Host, req.URL.Path)
		// log.Printf("--> Proxy. Mod incoming request header to %v", req.Header)

		// For debugging
		// dump, err := httputil.DumpRequest(req, true)
		// if err != nil {
		// 	return
		// }
		// log.Printf("--> Proxy. Before sending, dump req:\n%q", dump)
	}

	modifyResponse := func(res *http.Response) error {
		// log.Printf("<-- Proxy, modify response. Just got response, status: %v, response header: %v", res.Status, res.Header)
		log.Printf("<-- Just got response back from %v %v %v, status: %v", res.Request.URL.Scheme, res.Request.URL.Host, res.Request.URL.Path, res.Status)

		res.Header.Set("Access-Control-Allow-Origin", "*")

		for _, r := range serverConf.ResHeaderRules {
			r.Patch(originReqDomain2, originReqDomain3, res.Header)
		}

		for _, r := range serverConf.ResAddHeaderRules {
			r.Add(originReqDomain2, originReqDomain3, res.Header)
		}

		// extract gzipped content for modifications
		var resBody []byte
		switch res.Header.Get("Content-Encoding") {
		case "gzip":
			res.Header.Del("Content-Encoding")
			body, _ := gzip.NewReader(res.Body)
			resBody, _ = ioutil.ReadAll(body)
		default:
			resBody, _ = ioutil.ReadAll(res.Body)
		}

		modBody := bytes.NewReader(patchResponseBody(originReqHost, resBody, serverConf.BodyRules))
		res.Body = ioutil.NopCloser(modBody)
		res.Header["Content-Length"] = []string{fmt.Sprint(modBody.Len())}

		log.Printf("<<<----- Response sent.")

		return nil
	}

	// create the reverse proxy
	proxy := &httputil.ReverseProxy{Director: director, ModifyResponse: modifyResponse}
	log.Printf("----->>> Got original request %v %v from %v", req.Host, req.URL, req.RemoteAddr)

	proxy.ServeHTTP(res, req)
}

func disableTLSVerification() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func main() {
	configFilename := flag.String("conf", "config.json", "A config file in JSON format")
	flag.Parse()

	conf := config.Read(*configFilename)
	serverConf = configure.CompileConfig(conf)

	if conf.DisableTLS {
		disableTLSVerification()
	}

	server := http.NewServeMux()

	for _, static := range conf.StaticServers {
		fileServer := modifier.NoCache(http.FileServer(http.Dir(static.Root)))
		server.Handle(static.Prefix, http.StripPrefix(static.Prefix, fileServer))
	}

	server.HandleFunc("/", serveReverseProxy)

	log.Printf("Start a router at port %v", serverConf.ServerPort)
	if err := http.ListenAndServe(":"+serverConf.ServerPort, server); err != nil {
		panic(err)
	}

}
