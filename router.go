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
	"time"

	"github.com/andybalholm/brotli"

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

func serveReverseProxy(res http.ResponseWriter, req *http.Request) {
	originReqHost := modifier.CompileHostDomain(req.Host)

	director := func(req *http.Request) {
		log.Printf("----->>> Request received.")

		for _, r := range serverConf.SleepRules {
			if r.Match(req.URL.Path) && r.SleepSec > 0 {
				log.Printf("----- sleep for %v seconds", r.SleepSec)
				time.Sleep(time.Duration(r.SleepSec) * time.Second)
			}
		}

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
					directedHost = originReqHost.Host
				default:
					directedHost = dst.Host
				}

				switch dst.Server {
				case "":
					// use the server in the default rule
					directedServer = directedHost
				case "${HOST}":
					directedServer = originReqHost.Host
				default:
					directedServer = dst.Server
				}

				if dst.Path != "" {
					req.URL.Path = dst.Path
				}
				break
			}
		}

		for _, r := range serverConf.ReqAddCookieRules {
			if r.Match(req.URL.Path) {
				r.AddCookie(originReqHost, req)
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
			r.Patch(originReqHost, res.Header)
		}

		for _, r := range serverConf.ResAddHeaderRules {
			r.Add(originReqHost, res.Header)
		}

		// extract gzipped content for modifications
		var resBody []byte
		switch res.Header.Get("Content-Encoding") {
		case "gzip":
			res.Header.Del("Content-Encoding")
			body, _ := gzip.NewReader(res.Body)
			resBody, _ = ioutil.ReadAll(body)
		case "br":
			res.Header.Del("Content-Encoding")
			body := brotli.NewReader(res.Body)
			resBody, _ = ioutil.ReadAll(body)
		default:
			resBody, _ = ioutil.ReadAll(res.Body)
		}

		modBody := bytes.NewReader(patchResponseBody(originReqHost.Host, resBody, serverConf.BodyRules))
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

	isHTTPS := strings.Contains(strings.ToLower(serverConf.ServerScheme), "https")
	if isHTTPS {
		log.Printf("Start a router at port %v using https", serverConf.ServerPort)
		log.Fatal(http.ListenAndServeTLS(":"+serverConf.ServerPort, serverConf.ServerCert, serverConf.ServerKey, server))
	} else {
		log.Printf("Start a router at port %v using http", serverConf.ServerPort)
		log.Fatal(http.ListenAndServe(":"+serverConf.ServerPort, server))
	}
}
