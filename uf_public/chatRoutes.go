package uf_public

import (
	"github.com/gocql/gocql"
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

type ChatIdStruct struct {
	ChatId gocql.UUID
}

type SurferIdStruct struct {
	SurferId gocql.UUID
}

type PhoneStruct struct {
	Phone string
}

func HChatGetByPhone(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HChatGetByPhone called")
	in, out := PhoneStruct{}, ut.ChatPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	uf.Debug(gibs.ToddId)
	//TODO make admin only
	//rb.ExitIfPolicyFails(uf.PolicyToddIsSuperAdmin(gibs.ToddId))
	rb.ExitIfPolicyFails(uf.PolicyPublic())
	uf.Debug("in: ")
	uf.Debug(in.Phone)

	surfer, _ := uxt.SurferGetByPhone(&gibs, in.Phone)
	uf.Debug("authed")

	_ = uf.MakeRequest(
		&gibs,
		pile.Conf.AimAddress,
		uxt.AimChatGetBySurferIdV1,
		SurferIdStruct{SurferId: surfer.SurferId},
		&out,
	)
	uf.Debug("authed")

	msgs := ut.MsgsCollection{}
	_ = uf.MakeRequest(
		&gibs,
		pile.Conf.AimAddress,
		uxt.AimChatGetNewMsgsV1,
		in,
		&msgs,
	)
	out.Msgs = msgs.Collection
	uf.Debug("authed")

	rb.BuildResponse(ctx, &out)
}
