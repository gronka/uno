package uf_user

import (
	"github.com/gocql/gocql"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
)

type AddressIdStruct struct {
	AddressId gocql.UUID
}

var GetAddressError = uf.ApiError{
	Code: "get_address_error",
	Msg:  "get address error",
}

func HAddressGetById(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressGetById coming in")
	in, out := AddressIdStruct{}, AddressPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	out, err := LAddressGetById(&gibs, in.AddressId)

	if err != nil {
		uf.Debug("select address failed for AddressId: " + in.AddressId.String())
		uf.Debug(err.Error())
		rb.AddError(GetAddressError)
	}

	rb.BuildResponse(ctx, out)
}

func HAddressGetNonBuilderBySurferId(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressGetNonBuilderBySurferId coming in")
	in, out := SurferIdStruct{}, AddressPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Debug(in.SurferId)

	out, err := LAddressGetNonBuilderBySurferId(&gibs, in.SurferId)

	if err != nil {
		uf.Debug("select address failed for SurferId: " + in.SurferId.String())
		uf.Debug(err.Error())
		rb.AddError(GetAddressError)
	}

	rb.BuildResponse(ctx, out)
}

func HAddressGetBuilderBySurferId(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressGetBuilderBySurferId coming in")
	in, out := SurferIdStruct{}, AddressPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	out, err := LAddressGetBuilderForSurfer(&gibs, in.SurferId)

	if err != nil {
		uf.Debug("select address failed for SurferId: " + in.SurferId.String())
		uf.Debug(err.Error())
		rb.AddError(GetAddressError)
	}

	rb.BuildResponse(ctx, out)
}
