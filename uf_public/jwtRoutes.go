package uf_public

import (
	"github.com/gocql/gocql"
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
)

type MakeJwtIn struct {
	Uuid gocql.UUID
}

type MakeJwtOut struct {
	Jwt string
}

func HMakeJwt(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HMakeJwt called")
	in, out := MakeJwtIn{}, MakeJwtOut{}
	_, rb := uf.InitGibsForWeb(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyPublic())

	out.Jwt, _ = uf.CreateJwt(&pile.Conf, in.Uuid, "mock")

	rb.BuildResponse(ctx, out)
}
