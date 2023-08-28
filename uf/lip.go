package uf

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type HttpMethod string

const (
	HttpMethodGet  HttpMethod = "GET"
	HttpMethodPost            = "POST"
)

func makeUfUserUrl(path string) url.URL {
	return url.URL{
		Scheme:      "http",
		Opaque:      "",
		User:        &url.Userinfo{},
		Host:        "uf_user.com",
		Path:        path,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}
}

func CallUfUser(path string, out interface{}) {
	newUrl := makeUfUserUrl(path)

	jsonString, err := json.Marshal(out)
	if err != nil {
		panic("failed to marshal json")
	}

	http.Post(
		newUrl.String(),
		"application/json",
		bytes.NewBuffer(jsonString),
	)
}
