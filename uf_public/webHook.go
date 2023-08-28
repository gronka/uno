package uf_public

import (
	"github.com/gocql/gocql"
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_aim"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

type ChatSendMsgIn struct {
	ChatId         gocql.UUID
	SenderSurferId gocql.UUID
	Content        string
}

func HHookWeb(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HHookWeb")
	in, out := ChatSendMsgIn{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyToddIsSuperAdmin(gibs.ToddId))

	uf.Debug("1")
	surfer, ures := uxt.SurferGetById(&gibs, in.SenderSurferId)
	rb.AddErrors(ures.Errors)

	uf.Debug("1")
	pkg := uf_aim.MsgFromSurfer{
		SurferId:     surfer.SurferId,
		Recipient:    surfer.Phone,
		ChatId:       in.ChatId,
		SurferPhone:  surfer.Phone,
		ChatPlatform: ut.ChatPlatformWeb,
		Content:      in.Content,
		GroupId:      "",
	}
	uf.Debug("1")

	ures = uf.MakeRequest(
		&gibs,
		pile.Conf.AimAddress,
		uxt.AimMessageReceiveV1,
		pkg,
		uf.EmptyStruct{},
	)
	uf.Debug("1")

	rb.AddErrors(ures.Errors)

	uf.Debug("1")
	rb.BuildResponse(ctx, &out)
}

type HookWebCreateUserIn struct {
	Name      string
	Recipient string
}

func HHookWebCreateUser(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HHookWebCreateUser")
	in, out := HookWebCreateUserIn{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyPublic())

	uf.Trace("HHookWebCreateUser" + in.Recipient + ", for " + in.Name)

	//TODO: return error message if misformatted

	_, ures := uxt.SurferCreateFromPhoneOrEmail(
		&gibs,
		in.Name,
		in.Recipient,
	)

	rb.AddErrors(ures.Errors)

	uf.Debug("1")
	rb.BuildResponse(ctx, &out)
}
