package uf

import (
	"net/url"
)

func MakeUrl(domain string, path string, conf *Config) url.URL {
	return url.URL{
		Scheme:      conf.UfScheme,
		Opaque:      "",
		User:        nil,
		Host:        domain,
		Path:        path,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}
}
