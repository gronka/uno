package uf

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type ApiErrorsOut struct {
	Errors []ApiError
}

type ResponseBuilder struct {
	Body               interface{}
	Errors             []ApiError
	WasPolicyInspected bool
	//WereRequirementsInspected bool
}

func (rb *ResponseBuilder) AddError(apiError ApiError) {
	if rb.Errors == nil {
		rb.Errors = []ApiError{apiError}
	} else {
		rb.Errors = append(rb.Errors, apiError)
	}
}

func (rb *ResponseBuilder) AddErrors(apiErrors []ApiError) {
	if rb.Errors == nil {
		rb.Errors = apiErrors
	} else {
		rb.Errors = append(rb.Errors, apiErrors...)
	}
}

func (rb *ResponseBuilder) HasErrors() bool {
	if rb.Errors != nil && len(rb.Errors) > 0 {
		return true
	}
	return false
}

func (rb *ResponseBuilder) ExitIfPolicyFails(policyResult bool) {
	rb.WasPolicyInspected = true
	if !policyResult {
		panic(errors.New("policy:denied"))
	}
}

func (rb *ResponseBuilder) BuildResponse(
	ctx *fasthttp.RequestCtx,
	out interface{},
) {
	if !rb.WasPolicyInspected {
		panic(errors.New("policy:not_inspected"))
	}

	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.SetCanonical(
		[]byte("Content-Type"),
		[]byte("application/json"),
	)

	if rb.HasErrors() {
		rb.Body = ApiErrorsOut{rb.Errors}
	} else {
		rb.Body = out
	}

	if err := json.NewEncoder(ctx).Encode(rb.Body); err != nil {
		Check(err, "failed json encode for body")
	}
}
