package uf

import (
	"encoding/json"
	"fmt"

	"github.com/fasthttp/router"
	"github.com/gronka/tg"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func MakeRouter() *router.Router {
	r := router.New()
	r.PanicHandler = hPanicHandler

	r.GET("/", HHi)
	r.GET("/hello/{name}", hHello)
	r.NotFound = hNotFound

	return r
}

func HHi(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("HI")
}

func hHello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
	panic(errors.New("oh no"))
}

func hNotFound(ctx *fasthttp.RequestCtx) {
	tg.Warn("Route not found: " + string(ctx.Path()))
	ctx.Response.SetStatusCode(404)
}

func hPanicHandler(ctx *fasthttp.RequestCtx, err interface{}) {
	// Use %+v for detailed stack trace
	txt := fmt.Sprintf("%+v", err)
	tg.Error("API panic caught at router: " + txt)

	ctx.Response.Header.SetCanonical(
		[]byte("Content-Type"),
		[]byte("application/json"),
	)

	errCast, ok := err.(error)
	if ok {
		rootErr := errors.Cause(errCast)
		switch fmt.Sprintf("%v", rootErr) {
		case "gocql: no hosts available in the pool":
			connectCassandra(false)
			//TODO: email/call admin
			//TODO: driver automatically reconnects

		case "policy:denied":
			ctx.SetStatusCode(200)
			thisError := []ApiError{PolicyDeniedError}
			body := ApiErrorsOut{thisError}
			if err := json.NewEncoder(ctx).Encode(body); err != nil {
				panic(errors.Wrap(err, "failed json encode for body"))
			}

		case "policy:not_inspected":
			ctx.SetStatusCode(200)
			thisError := []ApiError{PolicyNotInspectedError}
			body := ApiErrorsOut{thisError}
			if err := json.NewEncoder(ctx).Encode(body); err != nil {
				panic(errors.Wrap(err, "failed json encode for body"))
			}

		case "requirements missing":
			ctx.SetStatusCode(200)
			thisError := []ApiError{RequirementsMissingError}
			body := ApiErrorsOut{thisError}
			if err := json.NewEncoder(ctx).Encode(body); err != nil {
				panic(errors.Wrap(err, "failed json encode for body"))
			}

		default:
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetBody([]byte(txt))
		}
	} else {
		panic("all errors should be wrapped with github.com/pkg/errors")
	}
}
