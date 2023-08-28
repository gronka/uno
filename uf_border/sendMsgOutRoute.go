package uf_border

import (
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
)

type SipOut struct {
	To          string `json:"to"`
	From        string `json:"from"`
	Body        string `json:"body"`
	MessageType string `json:"message_type"`
	NumRetries  int    `json:"num_retries"`
	SkipDefer   bool   `json:"skip_defer"`
}

type LoopMessageOut struct {
	// Recipient can be a number or email
	Recipient string `json:"recipient"`

	Text string `json:"text"`

	// MediaUrl max-length 256 characters
	MediaUrl string `json:"media_url"`

	// Sandbox enable for testing
	//Sandbox bool `json:"sandbox"`

	SenderName string `json:"sender_name"`

	// Timeout in seconds; minimum value is 5
	//Timeout int  `json:"timeout"`

	// Passthrough data for webhooks: 1000 char max
	Passthrough string `json:"passthrough"`

	// StatusCallback URL to receive message updates: 256 char max
	StatusCallback string `json:"status_callback"`

	// StatusCallbackHeader URL to receive message updates: 256 char max
	StatusCallbackHeader string `json:"status_callback_header"`
}

type LoopMessageResponse struct {
	Success string `json:"success"`
	// error codes: https://docs.loopmessage.com/messaging/send-message
	Code    string `json:"code"`
	Message string `json:"message"`
}

func HSendMsgOut(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HSendMsgOut")
	in, out := ut.MsgPg{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Debug(in.SenderSurferId.String())
	//surfer, _ := uxt.SurferGetById(&gibs, in.SenderSurferId)
	//chat, _ := uxt.ChatGetById(&gibs, in.ChatId)

	switch in.ChatPlatform {
	case ut.ChatPlatformSip:
		uf.Trace("sending sip message")
		sipOut := SipOut{
			To:          in.Recipient,
			From:        gibs.Conf.FridayyPhone,
			Body:        in.Content,
			MessageType: "SMS",
			NumRetries:  3,
			SkipDefer:   false,
		}

		_, err := SipRequest(
			gibs.Conf,
			gibs.Conf.SipUrl,
			uf.HttpMethodPost,
			sipOut,
			uf.EmptyStruct{},
		)
		if err != nil {
			uf.Glog(&gibs, uf.GlogStruct{
				Level:     uf.LevelError,
				Code:      "border.101",
				Interface: err.Error(),
			})
		}

	case ut.ChatPlatformLoop:
		uf.Trace("sending loop message")
		uf.Trace(in.Recipient)
		lmOut := LoopMessageOut{
			Recipient:  in.Recipient,
			Text:       in.Content,
			MediaUrl:   in.MediaUrl,
			SenderName: "fridayy@imsg.blue",
		}

		_, err := LoopMessageRequest(
			gibs.Conf,
			"https://server.loopmessage.com/api/v1/message/send/",
			uf.HttpMethodPost,
			lmOut,
			uf.EmptyStruct{},
		)
		if err != nil {
			uf.Glog(&gibs, uf.GlogStruct{
				Level:     uf.LevelError,
				Code:      "border.102",
				Interface: err.Error(),
			})
		}

	case ut.ChatPlatformWeb:
		uf.Trace("sending web message (do nothing)")

	}

	rb.BuildResponse(ctx, out)
}
