package uf_order

import (
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

// HOrderCreate creates the Order and a PaymentIntent to pay for it
func HPaymentAttemptCreate(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	in, out := ut.PaymentAttemptCreateIn{}, ut.PaymentAttemptCreateOut{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Trace("route: create payment attempt")

	order, _ := ut.LOrderSelectById(&gibs, in.OrderId)
	surfer, _ := uxt.SurferGetById(&gibs, in.SurferId)

	pm := GetDefaultPaymentMethod(&gibs, surfer.StripeCustomerId)

	uf.Debug("99999999999999999999999999999999999999999")
	uf.Debug(pm.ID)

	pa := PaymentAttemptPg{
		PaymentAttemptId: uf.RandomUUID(),
		OrderId:          order.OrderId,

		// Driver is stripe or btcpay, etc
		Driver:                "stripe",
		IsAutopaySelected:     false,
		Status:                "created",
		StripePaymentMethodId: pm.ID,

		Price:    order.TotalPrice(),
		Currency: order.Currency,

		TermsAcceptedVersion: "0.1",
		TimeCreated:          uf.NowStamp(),
		TimeUpdated:          uf.NowStamp(),
	}

	if in.UseSavedCard {
		uf.Debug("using saved card")

		pi, err := CreateStripePIWithSavedCard(
			&gibs,
			surfer.StripeCustomerId,
			order.TotalPrice(),
			order.Title,
			pm.ID,
		)
		//NOTE: set pa.ID
		uf.Debug("set payment intent id")
		pa.StripePaymentIntentId = pi.ID
		uf.Debug("set status")
		pa.Status = string(pi.Status)

		if err != nil {
			uf.Glog(&gibs, uf.GlogStruct{
				Level:     uf.LevelError,
				Code:      "order.111",
				Interface: err,
			})
		}

	} else {
		uf.Debug("not using saved card")

		sesh, err := CreateStripeSession(
			&gibs,
			surfer.StripeCustomerId,
			order.TotalPrice(),
			order.Title,
		)
		if err != nil {
			uf.Glog(&gibs, uf.GlogStruct{
				Level:     uf.LevelError,
				Code:      "order.111",
				Interface: err,
			})
		}
		//NOTE: set pa.ID
		uf.Debug("set url: " + sesh.URL)
		out.PaymentUrl = sesh.URL
		uf.Debug("set payment intent id")
		//uf.Debug(sesh)
		uf.Debug(sesh.PaymentIntent)
		//pa.StripePaymentIntentId = sesh
		//pa.StripePaymentIntentId = *sesh.PaymentIntent
		//pa.StripePaymentIntentId = string(sesh.PaymentIntent)

		sesh2, _ := GetCheckoutSession(&gibs, sesh.ID)
		uf.Debug("set payment intent id")
		uf.Debug(sesh2.PaymentIntent)
		pa.StripePaymentIntentId = sesh2.PaymentIntent.ID

	}

	uf.Debug("about to upsert")
	pa.LUpsert(&gibs)
	uf.Debug("upsert done")

	rb.BuildResponse(ctx, &out)
}
