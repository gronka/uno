package uf

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
)

type UfRequest struct {
	Url  string
	Body map[string]interface{}
}

type UfResponse struct {
	Resp    *http.Response
	Body    map[string]interface{}
	BodyObj interface{}
	UfId    gocql.UUID
	Errors  []ApiError
}

type AbstractBody struct {
	BodyObj interface{}
	Errors  []ApiError
}

type ApiErrorsFromBody struct {
	Errors []ApiError
}

func MakeRequest(
	gibs *Gibs,
	address string,
	path string,
	pkg interface{},
	out interface{},
) (ures UfResponse) {
	Trace("=====================================")

	url := MakeUrl(address, path, gibs.Conf)
	Trace("Internal request to " + url.String() + " with")
	pkgString, err := json.Marshal(pkg)
	if err != nil {
		ures.AddError(EncodeUfRequestBodyError)
	}

	Trace(pkg)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(pkgString))
	req.Header.Set("UfAuth", gibs.UfAuth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("UfId", gibs.UfId.String())
	req.Header.Set("UfKey", gibs.Conf.UfKey)

	ures.Resp, err = client.Do(req)
	if err != nil {
		ures.AddError(CustomError(err))
	}

	data, err := ioutil.ReadAll(ures.Resp.Body)
	if err != nil {
		panic(err.Error())
	}

	if ures.Resp == nil {
		ures.Body = make(map[string]interface{})
		ures.Errors = []ApiError{}
	} else {
		PrintResponseBody(data)

		reader := bytes.NewReader(data)
		err = json.NewDecoder(reader).Decode(&out)
		if err != nil {
			ures.AddError(DecodeUfResponseBodyError)
		}

		reader = bytes.NewReader(data)
		ures.unpackUfResponseErrors(reader)
		ures.BodyObj = out
	}

	Trace(ures.BodyObj)
	ures.LogErrors()

	return ures
}

func (ures *UfResponse) LogErrors() {
	if ures.Errors != nil {
		for _, ufError := range ures.Errors {
			Error(ufError)
		}
	}

	if ures.Resp == nil {
		Error("No response from server")
	}
	if ures.Resp.StatusCode < 200 ||
		ures.Resp.StatusCode > 299 {
		Error("Bad status code: " + strconv.Itoa(ures.Resp.StatusCode))
	}
}

func (ures *UfResponse) AddError(apiError ApiError) {
	if ures.Errors == nil {
		ures.Errors = []ApiError{apiError}
	} else {
		ures.Errors = append(ures.Errors, apiError)
	}
}

func (ures *UfResponse) AddErrors(apiErrors []ApiError) {
	if ures.Errors == nil {
		ures.Errors = apiErrors
	} else {
		for _, apiError := range apiErrors {
			ures.Errors = append(ures.Errors, apiError)
		}
	}
}

func (ures *UfResponse) Errored() bool {
	if ures.Errors != nil && len(ures.Errors) > 0 {
		return true
	}

	if ures.Resp == nil ||
		ures.Resp.StatusCode < 200 ||
		ures.Resp.StatusCode > 299 {
		return true
	}

	return false
}

func (ures *UfResponse) unpackUfResponseErrors(reader io.Reader) {
	respErrors := ApiErrorsFromBody{Errors: []ApiError{}}
	err := json.NewDecoder(reader).Decode(&respErrors)
	if err != nil {
		ures.AddError(DecodeUfResponseErrorsError)
	}

	for _, ApiError := range respErrors.Errors {
		ures.AddError(ApiError)
	}
}
