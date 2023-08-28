package uf_border

import (
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_aim"
	"gitlab.com/textfridayy/uno/uxt"
)

type JsonPhoneIn struct {
	Phone   string
	Content string
}

func (in *JsonPhoneIn) String() string {
	return in.Phone + in.Content
}

func HJsonPhoneReceive(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HJsonPhoneReceive called")
	in, out := JsonPhoneIn{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibsForSendBlue(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.AtLeastOnePolicyMustAllow([]bool{
		uf.PolicyCorrectSendBlueHookSecret(gibs),
		uf.PolicyCorrectSipHookSecret(gibs),
	}))

	tg.Debug(in.Phone)

	val := in.ValidateAllFields()
	if val.IsNotValid() {
		tg.Warn("JsonIn failed validation: " + in.String())
		rb.AddErrors(val.Errors)
		rb.BuildResponse(ctx, &out)
		return
	}

	surfer, _ := uxt.AuthViaSipHook(&gibs, pile, in.Phone)

	pkg2 := uf_aim.MsgFromSurfer{
		SurferId:     surfer.SurferId,
		SurferPhone:  surfer.Phone,
		ChatPlatform: "json",
		Content:      in.Content,
		GroupId:      "",
	}

	uf.Debug("about to make aim request for " + surfer.SurferId.String())

	uf.MakeRequest(
		&gibs,
		gibs.Conf.AimAddress,
		uxt.AimMessageReceiveV1,
		pkg2,
		uf.EmptyStruct{},
	)

	rb.BuildResponse(ctx, uf.EmptyStruct{})
}
