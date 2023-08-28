package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_public"
	"gitlab.com/textfridayy/uno/uxt"
)

func main() {
	fmt.Println("starting server")

	pile := uf.Pile{Conf: uf.Config{}}
	pile.InitFields("public")

	r := uf.MakeRouter()

	// health routes
	r.GET(uxt.PublicHealth, func(ctx *fasthttp.RequestCtx) {
		uf_public.HHealth(ctx, pile)
	})
	r.POST(uxt.PublicHealth, func(ctx *fasthttp.RequestCtx) {
		uf_public.HHealth(ctx, pile)
	})

	// hooks
	r.POST(uxt.PublicHookLoopReceive, func(ctx *fasthttp.RequestCtx) {
		uf_public.HLoopHook(ctx, pile)
	})
	r.GET(uxt.PublicHookLoopReceive, func(ctx *fasthttp.RequestCtx) {
		//uf_public.HLoopHook(ctx, pile)
		uf.HHi(ctx)
	})
	r.POST(uxt.PublicHookSipReceive, func(ctx *fasthttp.RequestCtx) {
		uf_public.HSipHook(ctx, pile)
	})

	// sign in routes
	r.POST(uxt.PublicToddSignInWithEmailV1, func(ctx *fasthttp.RequestCtx) {
		uf_public.HSignInWithEmail(ctx, pile)
	})
	r.POST(uxt.PublicToddSignInWithPhoneV1, func(ctx *fasthttp.RequestCtx) {
		uf_public.HSignInWithEmail(ctx, pile)
	})
	r.POST(uxt.PublicSurferGetByIdV1, func(ctx *fasthttp.RequestCtx) {
		uf_public.HSurferGetById(ctx, pile)
	})

	// chat routes
	r.POST(uxt.PublicChatGetByPhoneV1, func(ctx *fasthttp.RequestCtx) {
		uf_public.HChatGetByPhone(ctx, pile)
	})
	r.POST(uxt.PublicChatsGetListV1, func(ctx *fasthttp.RequestCtx) {
		uf_public.HChatsGetList(ctx, pile)
	})
	r.POST(uxt.PublicChatMsgsGetByIdV1, func(ctx *fasthttp.RequestCtx) {
		uf_public.HChatMsgsGetById(ctx, pile)
	})
	r.POST(uxt.PublicHookWebV1, func(ctx *fasthttp.RequestCtx) {
		uf_public.HHookWeb(ctx, pile)
	})
	r.POST(uxt.PublicHookWebCreateUserV1, func(ctx *fasthttp.RequestCtx) {
		uf_public.HHookWebCreateUser(ctx, pile)
	})

	uf.InitAndListenAndServe(pile.Conf, r)
}
