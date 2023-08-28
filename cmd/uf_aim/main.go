package main

import (
	"context"
	"fmt"

	"github.com/gronka/tg"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_aim"
	"gitlab.com/textfridayy/uno/uf_public"
	"gitlab.com/textfridayy/uno/uxt"
)

// var pgUrl = "postgres://postgres:postgres@uf_user_pg:5432/uf_user"
//var pgUrl = "postgres://postgres:postgres@localhost:5402/uf_user"

func main() {
	fmt.Println("starting server")

	pile := uf.Pile{Conf: uf.Config{}}
	pile.InitFields("aim")

	var err error = nil
	pile.Pool, err = pgxpool.Connect(context.Background(), pile.Conf.AimPgDsn)
	if err != nil {
		msg := "Unable to connect to " + pile.Conf.AimPgDsn
		tg.Error(msg)
	}
	defer pile.Pool.Close()

	r := uf.MakeRouter()

	// health routes
	r.GET(uxt.PublicHealth, func(ctx *fasthttp.RequestCtx) {
		uf_public.HHealth(ctx, pile)
	})
	r.POST(uxt.PublicHealth, func(ctx *fasthttp.RequestCtx) {
		uf_public.HHealth(ctx, pile)
	})

	// aim routes
	r.POST(uxt.AimMessageReceiveV1, func(ctx *fasthttp.RequestCtx) {
		uf_aim.HMessageReceive(ctx, pile)
	})

	// chat routes
	r.POST(uxt.AimChatGetBySurferIdV1, func(ctx *fasthttp.RequestCtx) {
		uf_aim.HChatGetBySurferId(ctx, pile)
	})
	r.POST(uxt.AimChatUpsertV1, func(ctx *fasthttp.RequestCtx) {
		uf_aim.HChatUpsert(ctx, pile)
	})
	r.POST(uxt.AimChatGetNewMsgsV1, func(ctx *fasthttp.RequestCtx) {
		uf_aim.HChatGetNewMsgs(ctx, pile)
	})
	r.POST(uxt.AimChatsGetListV1, func(ctx *fasthttp.RequestCtx) {
		uf_aim.HChatsGetList(ctx, pile)
	})

	uf.InitAndListenAndServe(pile.Conf, r)
}
