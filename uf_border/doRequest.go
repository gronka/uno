package uf_border

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/textfridayy/uno/uf"
)

func SipRequest(
	conf *uf.Config,
	reqUrl string,
	method uf.HttpMethod,
	in interface{},
	out interface{},
) (*http.Response, error) {

	parsed, _ := url.Parse(reqUrl)
	uf.Debug(parsed.String())
	bodyString, err := json.Marshal(in)

	req, err := http.NewRequest(
		string(method),
		parsed.String(),
		bytes.NewBuffer(bodyString))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sip request")
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return resp, errors.Wrap(err, "failed to do sip request")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, errors.New(fmt.Sprintf("Request to sip failed with Status Code %v", resp.StatusCode))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, errors.Wrap(err, "failed to decode sip request")
	}
	reader := bytes.NewReader(data)

	uf.Trace(uf.PrintResponseBody(data))

	err = json.NewDecoder(reader).Decode(&out)
	if err != nil {
		return resp, errors.Wrap(err, "failed to decode resp.Body: "+string(data))
	}

	return resp, nil
}

func LoopMessageRequest(
	conf *uf.Config,
	reqUrl string,
	method uf.HttpMethod,
	in interface{},
	out interface{},
) (*http.Response, error) {

	parsed, _ := url.Parse(reqUrl)
	uf.Debug(parsed.String())

	bodyString, err := json.Marshal(in)

	req, err := http.NewRequest(
		string(method),
		parsed.String(),
		bytes.NewBuffer(bodyString))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create LoopMessage request")
	}
	req.Header.Set("Authorization", conf.LoopAuthorization)
	req.Header.Set("Loop-Secret-Key", conf.LoopSecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return resp, errors.Wrap(err, "failed to do LoopMessage request")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, errors.New(fmt.Sprintf("Request to LoopMessage failed with Status Code %v", resp.StatusCode))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, errors.Wrap(err, "failed to decode LoopMessage request")
	}
	reader := bytes.NewReader(data)

	uf.Trace(uf.PrintResponseBody(data))

	err = json.NewDecoder(reader).Decode(&out)
	if err != nil {
		return resp, errors.Wrap(err, "failed to decode resp.Body: "+string(data))
	}

	return resp, nil
}
