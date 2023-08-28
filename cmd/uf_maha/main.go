package main

import (
	"context"
	"fmt"

	"github.com/gronka/tg"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_maha"
)

func main() {
	fmt.Println("starting server")

	pile := uf.Pile{Conf: uf.Config{}}
	pile.InitFields("maha")

	var err error = nil
	pile.Pool, err = pgxpool.Connect(context.Background(), pile.Conf.MahaPgDsn)
	if err != nil {
		msg := "Unable to connect to " + pile.Conf.MahaPgDsn
		tg.Error(msg)
	}
	defer pile.Pool.Close()

	r := uf.MakeRouter()

	r.POST(uf.MahaGlogV1, func(ctx *fasthttp.RequestCtx) {
		uf_maha.HGlogCreate(ctx, pile)
	})

	r.POST(uf.MahaGlogsGetByUfIdV1, func(ctx *fasthttp.RequestCtx) {
		uf_maha.HGlogsGetByUfId(ctx, pile)
	})

	uf.InitAndListenAndServe(pile.Conf, r)
}
