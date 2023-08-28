package uf_public

import (
	"github.com/gocql/gocql"
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

type GetChatMsgsIn struct {
	ChatId    gocql.UUID
	AfterTime int64
}

func HChatMsgsGetById(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HChatMsgsGetById called")
	in, out := ChatIdStruct{}, ut.MsgsCollection{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyToddIsSuperAdmin(gibs.ToddId))

	_ = uf.MakeRequest(
		&gibs,
		pile.Conf.AimAddress,
		uxt.AimChatGetNewMsgsV1,
		in,
		&out,
	)

	rb.BuildResponse(ctx, &out)
}
