package main

import (
	"context"
	"fmt"

	"github.com/gronka/tg"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
	"gitlab.com/textfridayy/uno/uxt"
)

// var pgUrl = "postgres://postgres:postgres@uf_user_pg:5432/uf_user"
//var pgUrl = "postgres://postgres:postgres@localhost:5402/uf_user"

func main() {
	fmt.Println("starting server")

	pile := uf.Pile{Conf: uf.Config{}}
	pile.InitFields("user")

	var err error = nil
	pile.Pool, err = pgxpool.Connect(context.Background(), pile.Conf.UserPgDsn)
	if err != nil {
		msg := "Unable to connect to " + pile.Conf.UserPgDsn
		tg.Error(msg)
	}
	defer pile.Pool.Close()

	r := uf.MakeRouter()

	r.POST(uxt.UserSurferCreateStripeCustomerIdV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HSurferCreateStripeCustomerId(ctx, pile)
	})

	r.POST(uxt.UserSurferGetOrCreateFromPhoneOrEmailV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HSurferGetOrCreateFromPhoneOrEmail(ctx, pile)
	})

	r.POST(uxt.UserSurferGetByEmailV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HSurferGetByEmail(ctx, pile)
	})
	r.POST(uxt.UserSurferGetByPhoneV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HSurferGetByPhone(ctx, pile)
	})
	r.POST(uxt.UserSurferGetByIdV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HSurferGetBySurferId(ctx, pile)
	})

	r.POST(uxt.UserAddressCreateBuilderV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressCreateBuilder(ctx, pile)
	})
	r.POST(uxt.UserAddressDeleteV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressDelete(ctx, pile)
	})
	r.POST(uxt.UserAddressDeleteBuilderV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressDeleteBuilder(ctx, pile)
	})
	r.POST(uxt.UserAddressGetBuilderBySurferIdV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressGetBuilderBySurferId(ctx, pile)
	})
	r.POST(uxt.UserAddressGetNonBuilderBySurferIdV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressGetNonBuilderBySurferId(ctx, pile)
	})
	r.POST(uxt.UserAddressGetByIdV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressGetById(ctx, pile)
	})

	r.POST(uxt.UserAddressUpdateNameV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressUpdateName(ctx, pile)
	})
	r.POST(uxt.UserAddressUpdatePostalPlusV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressUpdatePostalPlus(ctx, pile)
	})
	r.POST(uxt.UserAddressUpdateLine1V1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressUpdateAddressLine1(ctx, pile)
	})
	r.POST(uxt.UserAddressUpdateLine2V1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressUpdateAddressLine2(ctx, pile)
	})
	r.POST(uxt.UserAddressValidateUspsV1, func(ctx *fasthttp.RequestCtx) {
		uf_user.HAddressValidateUsps(ctx, pile)
	})

	uf.InitAndListenAndServe(pile.Conf, r)
}
