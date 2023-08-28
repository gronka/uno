package uf_public

import (
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

func HChatsGetList(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HChatGetByPhone called")
	in, out := uf.EmptyStruct{}, ut.ChatCollection{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	uf.Debug(gibs.ToddId)
	rb.ExitIfPolicyFails(uf.PolicyToddIsSuperAdmin(gibs.ToddId))

	_ = uf.MakeRequest(
		&gibs,
		pile.Conf.AimAddress,
		uxt.AimChatsGetListV1,
		uf.EmptyStruct{},
		&out,
	)
	uf.Debug("authed")

	rb.BuildResponse(ctx, out)
}
