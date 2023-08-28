package uf_order

import (
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
	"gitlab.com/textfridayy/uno/ut"
)

func HOrdersGetBySurferId(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	in, out := uf_user.SurferIdStruct{}, OrderPgCollection{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Trace("route: get all orders")

	var err error
	out, err = LOrdersSelectBySurferId(&gibs, in.SurferId)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "order.112",
			Interface: err,
		})
	}

	rb.BuildResponse(ctx, &out)
}

func HOrderGetByIdDetails(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	in, out := OrderIdStruct{}, ut.OrderDetails{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Trace("route: get order details")
	order, err := ut.LOrderSelectById(&gibs, in.OrderId)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "order.113",
			Interface: err,
		})
	}

	out, err = ut.GetOrderDetailsByDriver(&gibs, &order)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "order.114",
			Interface: err,
		})
	}
	uf.Trace(out)

	rb.BuildResponse(ctx, &out)
}

func HOrderCancel(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	in, out := OrderIdStruct{}, ut.CancelDetails{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Trace("route: order cancel")
	order, err := ut.LOrderSelectById(&gibs, in.OrderId)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "order.115",
			Interface: err,
		})
	}

	out, err = ut.OrderCancelByDriver(&gibs, &order)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "order.116",
			Interface: err,
		})
	}
	uf.Trace(out)

	rb.BuildResponse(ctx, &out)
}

func HOrderGetBySurferIdNewest(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	in, out := OrderIdStruct{}, ut.OrderPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Trace("route: order get newest")
	out, err := ut.LOrderSelectById(&gibs, in.OrderId)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "order.117",
			Interface: err,
		})
	}

	rb.BuildResponse(ctx, &out)
}
