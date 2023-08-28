package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_border"
	"gitlab.com/textfridayy/uno/uxt"
)

func main() {
	fmt.Println("starting server")

	pile := uf.Pile{Conf: uf.Config{}}
	pile.InitFields("border")

	r := uf.MakeRouter()

	r.POST(uxt.BorderJsonHookReceiveV1, func(ctx *fasthttp.RequestCtx) {
		uf_border.HJsonPhoneReceive(ctx, pile)
	})

	r.POST(uxt.BorderSendMsgOutV1, func(ctx *fasthttp.RequestCtx) {
		uf_border.HSendMsgOut(ctx, pile)
	})

	uf.InitAndListenAndServe(pile.Conf, r)
}
