package uf_public

import (
	"strings"

	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

type LoopHookIn struct {
	MesssageId string `json:"messsage_id"`

	// WebhookId is event identifier
	WebhookId string `json:"webhook_id"`

	// AlertType values: message_scheduled, message_failed, message_sent,
	// message_reply, and message_reaction
	AlertType string `json:"alert_type"`

	// Recipient for the hook is actually the sender can be a number or email.
	Recipient string `json:"recipient"`

	Text    string `json:"text"`
	Success bool   `json:"success"`

	// Sandbox: optional - is true if this is a reply coming from a sandbox action
	Sandbox bool `json:"sandbox"`

	// Reaction is only set for alert_type message_reaction. Possible values:
	// love, like, dislike, laugh, exclaim, question, unknown
	Reaction string `json:"reaction"`

	// ErrorCode is only set for alert_type message_failed
	ErrorCode int `json:"error_code"`

	// Passthrough data for webhooks: 1000 char max
	Passthrough string `json:"passthrough"`

	// IsPhone is not from LoopMessage
	IsPhone bool
}

type Simp struct {
	Ok bool
}

func HLoopHook(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HLoopHook called")
	in, out := LoopHookIn{}, Simp{}
	gibs, rb := uf.InitGibsForLoop(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyCorrectLoopHookAuthorization(gibs))

	uf.Debug("1")

	uf.Trace("phoneOrEmail: " + in.Recipient)
	if strings.Contains(in.Recipient, "@") {
		in.IsPhone = false
	} else {
		in.IsPhone = true
	}

	surfer, _ := uxt.AuthViaLoopHook(&gibs, pile, in.IsPhone, in.Recipient)
	uf.Debug("2")
	chat, ures := uxt.ChatGetBySurferId(&gibs, surfer.SurferId)
	uf.Debug("3")

	SendMsgToAim(
		ut.ChatPlatformLoop,
		&chat,
		&in,
		&gibs,
		&surfer,
		&rb,
		&ures)

	uf.Debug("HLoopHook DONE")
	out.Ok = true
	rb.BuildResponse(ctx, &out)
}
