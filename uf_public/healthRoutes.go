package uf_public

import (
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
)

func HHealth(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HHealth called")
	in, out := uf.EmptyStruct{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibsForWeb(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyPublic())

	uf.Glog(&gibs, uf.GlogStruct{
		Code:  "099",
		Level: uf.LevelDebug,
	})

	rb.BuildResponse(ctx, out)
}
