// Package configure translates configs into rules and settings for the router server.
package configure

import (
	"net/url"
	"regexp"

	"github.com/bradleyzhou/dev-router/config"
	"github.com/bradleyzhou/dev-router/modifier"
)

// Config contains the rules and settings used by the router.
type Config struct {
	ServerPort   string
	ServerScheme string

	BodyRules         []modifier.PatchBodyRule
	ResHeaderRules    []modifier.PatchHeaderRule
	ResAddHeaderRules []modifier.AddHeaderRule
	ReqAddCookieRules []modifier.AddRequestCookieRule

	DefaultRequestDispatchRule modifier.RequestDispatchRule
	RequestDispatchRules       []modifier.RequestDispatchRule
	SleepRules                 []modifier.RequestSleepRule
}

// CompileConfig turns a config (read from file) into modifier rules and settings used by the router.
func CompileConfig(conf config.Config) Config {
	return Config{
		ServerPort:                 conf.Server.Port,
		ServerScheme:               conf.Server.Scheme,
		BodyRules:                  compileBodyPatchers(conf),
		DefaultRequestDispatchRule: compileDefaultDispatcher(conf.Upstream),
		RequestDispatchRules:       compileDispatchers(conf),
		ResHeaderRules:             compileHeaderPatchers(conf.ResponseModifiers.HeaderPatchers),
		ResAddHeaderRules:          compileHeaderAdders(conf.ResponseHeaderAdders),
		ReqAddCookieRules:          compileReqCookieAdders(conf.RequestCookieAdders),
		SleepRules:                 compileSleepers(conf.RequestSleepers),
	}
}

func compileBodyPatchers(conf config.Config) []modifier.PatchBodyRule {
	raw := conf.ResponseModifiers.BodyPatchers
	compiled := make([]modifier.PatchBodyRule, len(raw))
	for i, r := range raw {
		compiled[i].Matcher = regexp.MustCompile(r.Matcher)
		compiled[i].Replacer = []byte(r.Replacer)
	}
	return compiled
}

func compileReqCookieAdders(raw []config.AddRequestCookieRule) []modifier.AddRequestCookieRule {
	compiled := make([]modifier.AddRequestCookieRule, len(raw))
	for i, r := range raw {
		compiled[i].PathMatcher = regexp.MustCompile(r.PathMatcher)
		compiled[i].CookieAdder = modifier.CookieAdder{
			Name:  r.Name,
			Value: r.Value,
		}
	}
	return compiled
}

func compileHeaderAdders(raw []config.AddHeaderRule) []modifier.AddHeaderRule {
	compiled := make([]modifier.AddHeaderRule, len(raw))
	for i, r := range raw {
		compiled[i].Name = r.Name
		compiled[i].Value = r.Value
	}
	return compiled
}

func compileHeaderPatchers(raw []config.PatchHeaderRule) []modifier.PatchHeaderRule {
	compiled := make([]modifier.PatchHeaderRule, len(raw))
	for i, r := range raw {
		compiled[i].Name = r.Name
		compiled[i].Matcher = regexp.MustCompile(r.Matcher)
		compiled[i].Replacer = r.Replacer
	}
	return compiled
}

func compileSleepers(raw []config.SleepRule) []modifier.RequestSleepRule {
	compiled := make([]modifier.RequestSleepRule, len(raw))
	for i, r := range raw {
		compiled[i].PathMatcher = regexp.MustCompile(r.PathMatcher)
		compiled[i].SleepSec = r.SleepSeconds
	}
	return compiled
}

func compileDispatchers(conf config.Config) []modifier.RequestDispatchRule {
	raw := conf.RequestDispatchers
	compiled := make([]modifier.RequestDispatchRule, len(raw))
	for i, r := range raw {
		compiled[i].DstServer = r.Destination.Server
		compiled[i].DstScheme = r.Destination.Scheme
		compiled[i].DstHost = r.Destination.Host
		compiled[i].PathMatcher = regexp.MustCompile(r.PathMatcher)
		compiled[i].PathReplacer = r.Destination.PathReplacer
	}
	return compiled
}

func compileDefaultDispatcher(URL string) modifier.RequestDispatchRule {
	if target, err := url.Parse(URL); err == nil {
		return modifier.RequestDispatchRule{
			DstScheme: target.Scheme,
			DstHost:   target.Host,
		}
	}
	return modifier.RequestDispatchRule{}
}
