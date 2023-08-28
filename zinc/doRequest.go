package zinc

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

func DoRequest(
	conf *uf.Config,
	reqUrl string,
	method uf.HttpMethod,
	in interface{},
	params map[string]string,
	out interface{},
) (*http.Response, error) {

	if in == nil {
		in = make(map[string]string)
	}

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	/* another way
	values := url.Values{
	"page":     []string{"1"},
	"query":    []string{query},
	"retailer": []string{retailer},
	}
	*/

	parsed, _ := url.Parse(reqUrl)
	parsed.RawQuery = values.Encode()
	parsed.User = url.User(conf.ZincToken)
	uf.Debug(parsed.String())

	var bodyBytes bytes.Buffer
	if in != nil {
		bodyString, err := json.Marshal(in)
		uf.Trace(string(bodyString))
		uf.Trace(parsed.String())
		if err != nil {
			return nil, errors.Wrap(err, "Failed to stringify body")
		}
		bodyBytes := bytes.NewBuffer(bodyString)
		uf.Trace(1)
		uf.Trace(bodyBytes)
	}

	req, err := http.NewRequest(string(method), parsed.String(), &bodyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create zinc request")
	}

	uf.Trace(2)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return resp, errors.Wrap(err, "failed to do zinc request")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, errors.New(fmt.Sprintf("Request to Zinc failed with Status Code %v", resp.StatusCode))
	}

	uf.Trace(3)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, errors.Wrap(err, "failed to decode zinc request")
	}
	reader := bytes.NewReader(data)

	uf.Trace(uf.PrintResponseBody(data))

	err = json.NewDecoder(reader).Decode(&out)
	if err != nil {
		return resp, errors.Wrap(err, "failed to decode resp.Body: "+string(data))
	}

	return resp, nil
}
