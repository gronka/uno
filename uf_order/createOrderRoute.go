package uf_order

import (
	"errors"

	"github.com/gocql/gocql"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

var CreateBasketError = uf.ApiError{
	Code: "create_basket_error",
	Msg:  "Create basket error.",
}

// HOrderCreate creates the Order and a PaymentIntent to pay for it
func HOrderCreate(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	in, out := ut.OrderCreateIn{}, ut.OrderCreateOut{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Trace("create order")

	surfer, _ := uxt.SurferGetById(&gibs, in.SurferId)

	if surfer.StripeCustomerId == "" {
		uxt.SurferCreateStripeCustomerId(&gibs, surfer.SurferId)
		surfer, _ = uxt.SurferGetById(&gibs, in.SurferId)
	}

	// create basket items
	orderId := uf.RandomUUID()
	out.OrderId = orderId
	basketItemIds := []gocql.UUID{}
	basketPrice := 0
	for _, item := range in.BasketItems {
		item.BasketItemId = uf.RandomUUID()
		item.OrderId = orderId
		err := item.LUpsert(&gibs)
		if err != nil {
			rb.AddError(uf.CustomError(
				errors.New("create_basket_error:" + item.ProductId.String())))
		}
		basketItemIds = append(basketItemIds, item.BasketItemId)
		basketPrice += item.UnitPrice * item.Quantity
	}

	address, _ := uxt.AddressGetById(&gibs, in.ShippingAddressId)
	tax := pile.TaxTable.GetRate(address.Postal)

	now := uf.NowStamp()
	order := ut.OrderPg{
		OrderId:              orderId,
		SurferId:             surfer.SurferId,
		SurferSubscriptionId: 0,
		BasketItemIds:        basketItemIds,

		Status: "created",

		Title:         in.Title,
		BasketPrice:   basketPrice,
		Tax:           tax,
		Margin:        uf.AddOurMargin(basketPrice),
		Currency:      "USD",
		ShippingDays:  in.ShippingDays,
		ShippingPrice: in.ShippingPrice,

		ShippingAddressId: address.AddressId,
		ShippingName:      address.Name,
		ShippingLine1:     address.AddressLine1,
		ShippingLine2:     address.AddressLine2,
		ShippingCity:      address.City,
		ShippingState:     address.State,
		ShippingPostal:    address.Postal,

		// OrderSystem would be stripe, btcpay, etc
		Driver:      "stripe",
		TimeCreated: now,
		TimeUpdated: now,
	}

	err := order.LUpsert(&gibs)

	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "order.110",
			Interface: err,
		})
	}

	pm := GetDefaultPaymentMethod(&gibs, surfer.StripeCustomerId)
	out.PaymentMethodId = pm.ID
	out.IsPaymentMethodExpired = IsPaymentMethodExpired(pm)

	if pm.Card != nil {
		out.CardLast4 = pm.Card.Last4
	}

	uf.Warn(out)

	//TODO get surfer subscription

	uf.Debug("orderId" + orderId.String())
	rb.BuildResponse(ctx, &out)
}
