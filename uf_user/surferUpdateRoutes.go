package uf_user

import (
	"github.com/gocql/gocql"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/customer"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
)

type SurferUpdateNameStruct struct {
	SurferId gocql.UUID
	Name     string
}

func HSurferUpdateName(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HSurferUpdateName coming in")
	in, out := SurferUpdateNameStruct{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Debug(in)

	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `UPDATE surfers SET
		name = $1,
		WHERE surfer_id = $2`,
		in.Name,
		in.SurferId,
	)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "user.user.101",
			Interface: err,
		})
	}

	rb.BuildResponse(ctx, out)
}

func HSurferCreateStripeCustomerId(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HSurferCreateStripeCustomerId coming in")
	in, out := SurferIdStruct{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	surfer, err := LSurferGetById(&gibs, in.SurferId)
	if err != nil {
		rb.AddError(GetUserError)
	}
	if surfer.StripeCustomerId != "" {
		rb.BuildResponse(ctx, out)
		return
	}

	stripe.Key = gibs.Conf.StripeSecretKey
	params := &stripe.CustomerParams{
		Description:      stripe.String("fridayy customer"),
		Email:            stripe.String(surfer.Email),
		Phone:            stripe.String(surfer.Phone),
		PreferredLocales: stripe.StringSlice([]string{"en", "es"}),
	}

	cust, err := customer.New(params)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "user.user.102",
			Interface: err,
		})
	}

	_, err = gibs.Pile.Pool.Exec(gibs.Ctx, `UPDATE surfers SET
		stripe_customer_id = $1
		WHERE surfer_id = $2`,
		cust.ID,
		in.SurferId,
	)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "user.user.103",
			Interface: err,
		})
	}

	rb.BuildResponse(ctx, out)
}
