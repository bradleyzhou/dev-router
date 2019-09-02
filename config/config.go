// Package config represents the user provided config file structure. Right now with JSON support.
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Server TODO
type Server struct {
	Port   string `json:"port"`
	Scheme string `json:"scheme"`
}

// StaticServer TODO
type StaticServer struct {
	Prefix string `json:"prefix"`
	Root   string `json:"root"`
}

// DispatchDestination is the destination part of a request dispatcher.
type DispatchDestination struct {
	Server       string `json:"server"`
	Scheme       string `json:"scheme"`
	Host         string `json:"host"`
	PathReplacer string `json:"path"`
	URL          string `json:"url"`
}

// DispatchRule TODO
type DispatchRule struct {
	PathMatcher string              `json:"path"`
	Destination DispatchDestination `json:"destination"`
}

// ResponseModifier TODO
type ResponseModifier struct {
	HeaderPatchers []PatchHeaderRule `json:"header"`
	BodyPatchers   []PatchBodyRule   `json:"body"`
}

// PatchBodyRule TODO
type PatchBodyRule struct {
	Matcher  string `json:"matcher"`
	Replacer string `json:"replacer"`
}

// PatchHeaderRule TODO
type PatchHeaderRule struct {
	Name     string `json:"name"`
	Matcher  string `json:"matcher"`
	Replacer string `json:"replacer"`
}

// AddHeaderRule TODO
type AddHeaderRule struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AddRequestCookieRule TODO
type AddRequestCookieRule struct {
	PathMatcher string `json:"path"`
	Name        string `json:"name"`
	Value       string `json:"value"`
}

// Config TODO
type Config struct {
	Server        Server         `json:"server"`
	Upstream      string         `json:"upstream"`
	DisableTLS    bool           `json:"disableTLSVerify"`
	StaticServers []StaticServer `json:"staticServers"`

	RequestDispatchers   []DispatchRule         `json:"requestDispatchers"`
	ResponseModifiers    ResponseModifier       `json:"responseModifiers"`
	ResponseHeaderAdders []AddHeaderRule        `json:"addResponseHeader"`
	RequestCookieAdders  []AddRequestCookieRule `json:"addRequestCookie"`
}

// Read reads the server configs from a .json file
func Read(filename string) Config {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	content, _ := ioutil.ReadAll(f)

	var config Config
	json.Unmarshal(content, &config)
	return config
}
