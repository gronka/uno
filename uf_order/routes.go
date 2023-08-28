package uf_order

import (
	"github.com/gocql/gocql"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/scaffold"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uxt"
	"gitlab.com/textfridayy/uno/zinc"
)

var OrderCreateError = uf.ApiError{
	Code: "order_create_error",
	Msg:  "Error creating order for surfer",
}

type CreateOrderIn struct {
	SurferId gocql.UUID
	OrderId  gocql.UUID
	Testing  bool
}

func HOrderCreateOld(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	in := CreateOrderIn{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	surfer, _ := uxt.SurferGetById(&gibs, in.SurferId)
	shippingAddress, _ := uxt.AddressGetNonBuilderBySurferId(&gibs, in.SurferId)
	shippingAddressObject := shippingAddress.AsZincAddressObject()

	//TODO: pass query instead of this junk
	shippingInfo := zinc.ShippingObject{
		OrderBy:  "price",
		MaxPrice: 0,
	}

	//TODO
	products := []zinc.ProductObject{scaffold.ProductObject01()}

	_, _, _, err := ExecuteCreateOrder(
		&pile.Conf,
		surfer,
		shippingAddressObject,
		nil, // billingAddress
		shippingInfo,
		products,
		0, // maxPrice
		in.Testing,
	)
	if err != nil {
		uf.Warn(err)
	}

	uf.Debug("return message sending")
	ctx.Response.SetStatusCode(200)
}
