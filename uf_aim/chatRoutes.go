package uf_aim

import (
	"github.com/gocql/gocql"
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
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

func HChatGetBySurferId(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HChatGetBySurferId called")
	in := SurferIdStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	chat, err := ut.LChatGetBySurferId(&gibs, in.SurferId)
	if err != nil {
		uf.Error(err)
	}

	msgs, err := ut.LChatGetMsgs(&gibs, chat.ChatId)
	chat.Msgs = msgs

	rb.BuildResponse(ctx, &chat)
}

func HChatGetNewMsgs(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HChatGetNewMsgs called")
	in, out := ChatIdStruct{}, ut.MsgsCollection{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	var err error
	out.Collection, err = ut.LChatGetMsgs(&gibs, in.ChatId)
	if err != nil {
		uf.Error(err)
	}
	rb.BuildResponse(ctx, &out)
}

func HChatsGetList(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HChatGetList called")
	in, out := uf.EmptyStruct{}, ut.ChatCollection{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	var err error
	out.Collection, err = ut.LChatsGetList(&gibs)
	if err != nil {
		uf.Error(err)
	}
	rb.BuildResponse(ctx, &out)
}

func HChatUpsert(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HChatUpsert called")
	in := ut.ChatPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	chat := in.LUpsert(&gibs)
	rb.BuildResponse(ctx, &chat)
}
