package ut

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/textfridayy/uno/uf"
)

type SmsData struct {
	Phone    string
	Content  string
	MediaUrl string
}

type SmsDataForSendBlue struct {
	Number   string `json:"number"`
	Content  string `json:"content"`
	MediaUrl string `json:"media_url"`
}

func (out *SmsDataForSendBlue) Send(conf uf.Config) error {
	endpoint := url.URL{
		Scheme:      "https",
		Opaque:      "",
		User:        nil,
		Host:        "api.sendblue.co",
		Path:        "api/send-message",
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}
	uf.Debug(endpoint.String())
	end := "https://api.sendblue.co/api/send-message"

	out.Content = "shorter"
	pkgString, err := json.Marshal(out)
	if err != nil {
		return errors.New("failed to encode")
	}

	uf.Debug(bytes.NewBuffer(pkgString))

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		//endpoint.String(),
		end,
		bytes.NewBuffer(pkgString),
	)
	if err != nil {
		return errors.Wrap(err, "failed to build sendblue request")
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("sb-api-key-id", conf.SendBlueApiKeyId)
	req.Header.Set("sb-api-secret-key", conf.SendBlueApiSecret)

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send sendblue request")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(
			fmt.Sprintf("sendblue request returned error code: %v", resp.StatusCode))
	}

	amap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&amap)
	if err != nil {
		return err
	}
	uf.Debug(amap)

	return nil
}
