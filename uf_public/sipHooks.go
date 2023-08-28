package uf_public

import (
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

type SipHookIn struct {
	Type              string `json:"type"`
	DestinationNumber string `json:"destination_number"`
	CallerIdNumber    string `json:"caller_id_number"`
	Text              string `json:"text"`
}

func HSipHook(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HSipHook called")
	in, out := SipHookIn{}, Simp{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	uf.Trace(gibs.SipHookSecret)
	uf.Trace(gibs.Conf.SipHookSecret)
	rb.ExitIfPolicyFails(uf.PolicyCorrectSipHookSecret(gibs))

	uf.Debug("1")

	surfer, _ := uxt.AuthViaSipHook(&gibs, pile, in.CallerIdNumber)
	uf.Debug("2")
	chat, ures := uxt.ChatGetBySurferId(&gibs, surfer.SurferId)
	uf.Debug("3")

	inAsLoopIn := LoopHookIn{
		Recipient: in.CallerIdNumber,
		Text:      in.Text,
		IsPhone:   true,
	}

	SendMsgToAim(
		ut.ChatPlatformSip,
		&chat,
		&inAsLoopIn,
		&gibs,
		&surfer,
		&rb,
		&ures)

	uf.Debug("HSipHook DONE")
	out.Ok = true
	rb.BuildResponse(ctx, &out)
}
