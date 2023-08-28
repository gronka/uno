package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gronka/tg"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_order"
	"gitlab.com/textfridayy/uno/uxt"
)

func main() {
	fmt.Println("starting server")

	pile := uf.Pile{Conf: uf.Config{}}
	pile.InitFields("order")

	var err error = nil
	pile.Pool, err = pgxpool.Connect(context.Background(), pile.Conf.OrderPgDsn)
	if err != nil {
		msg := "Unable to connect to " + pile.Conf.OrderPgDsn
		tg.Error(msg)
	}
	defer pile.Pool.Close()

	r := uf.MakeRouter()

	r.POST(uxt.OrderPaymentAttemptCreateV1, func(ctx *fasthttp.RequestCtx) {
		uf_order.HPaymentAttemptCreate(ctx, pile)
	})

	r.POST(uxt.OrderOrderCancelV1, func(ctx *fasthttp.RequestCtx) {
		uf_order.HOrderCancel(ctx, pile)
	})
	r.POST(uxt.OrderOrderCreateV1, func(ctx *fasthttp.RequestCtx) {
		uf_order.HOrderCreate(ctx, pile)
	})
	r.POST(uxt.OrderOrderGetByIdDetailsV1, func(ctx *fasthttp.RequestCtx) {
		uf_order.HOrderCreate(ctx, pile)
	})
	r.POST(uxt.OrderOrderGetBySurferIdNewestV1, func(ctx *fasthttp.RequestCtx) {
		uf_order.HOrderGetBySurferIdNewest(ctx, pile)
	})

	r.POST(uxt.OrderOrdersGetBySurferIdV1, func(ctx *fasthttp.RequestCtx) {
		uf_order.HOrdersGetBySurferId(ctx, pile)
	})

	uf.InitAndListenAndServe(pile.Conf, r)
}

func sleepyPrint() {
	for true {
		time.Sleep(3 * time.Second)
		uf.Trace("3 seconds")
	}
}
