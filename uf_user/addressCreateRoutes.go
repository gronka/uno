package uf_user

import (
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
)

var BuildAddressError = uf.ApiError{
	Code: "build_address_error",
	Msg:  "build address error",
}

func HAddressCreateBuilder(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressCreateBuilder coming in")
	in, out := SurferIdStruct{}, AddressPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	now := uf.NowStamp()
	builder := AddressPg{
		AddressId:   uf.RandomUUID(),
		SurferId:    in.SurferId,
		IsBuilder:   true,
		TimeCreated: now,
		TimeUpdated: now,
	}
	err := builder.LAddressUpsert(&gibs)
	//err := LAddressCreateBuilder(&gibs, in.SurferId)

	if err != nil {
		uf.Debug("create builder failed for SurferId: " + in.SurferId.String())
		rb.AddError(uf.CustomError(err))
		panic(err)
	}

	rb.BuildResponse(ctx, out)
}

func HAddressDelete(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressDelete coming in")
	in, out := AddressIdStruct{}, AddressPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `DELETE FROM addresses 
		WHERE address_id = $1`,
		in.AddressId,
	)

	if err != nil {
		uf.Debug("Delete address failed for AddressId: " + in.AddressId.String())
		rb.AddError(uf.CustomError(err))
	}

	rb.BuildResponse(ctx, out)
}

func HAddressDeleteBuilder(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressDeleteBuilder coming in")
	in, out := SurferIdStruct{}, AddressPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `DELETE FROM addresses 
		WHERE surfer_id = $1 AND is_builder = $2`,
		in.SurferId,
		true,
	)

	if err != nil {
		uf.Debug("Delete builder failed for SurferId: " + in.SurferId.String())
		rb.AddError(uf.CustomError(err))
	}

	rb.BuildResponse(ctx, out)
}
