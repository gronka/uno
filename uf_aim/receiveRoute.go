package uf_aim

import (
	"github.com/gocql/gocql"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

var AimCreateError = uf.ApiError{
	Code: "aim_create_error",
	Msg:  "Error creating aim for surfer",
}

var AimUpdateError = uf.ApiError{
	Code: "aim_update_error",
	Msg:  "Error updating aim for surfer.",
}

var GetChatError = uf.ApiError{
	Code: "get_chat_error",
	Msg:  "Error getting chat.",
}

var SaveMessageError = uf.ApiError{
	Code: "save_message_error",
	Msg:  "Error saving message.",
}

type MsgFromSurfer struct {
	// SMS, mobile app, etc
	SurferId    gocql.UUID
	SurferPhone string
	SurferEmail string
	Recipient   string
	ChatId      gocql.UUID

	// ChatPlatform can be loop_message or sip
	ChatPlatform ut.ChatPlatform
	// Message content
	Content string

	GroupId string
}

/*
steps to determine aim:
 1. read aim from DB (aim might be predetermined)
    1.a) if aim.content can be acted upon, perform the action
    1.b) if aim.content cannot be acted upon, then ask friday the user's aim
 2. update aim_series according to the action taken, so that on the next request

we will know where the user is coming from and where they are going.
*/
func HMessageReceive(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	in, out := MsgFromSurfer{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))
	uf.Debug(in)

	uf.Debug("get aim by surfer id")
	aimInfo, err := LAimGetBySurferId(&gibs, in.SurferId)
	if err != nil {
		uf.Warn(err)
		uf.Debug("create aim")
		err = LAimCreate(&gibs, in.SurferId)
		if err != nil {
			uf.Error(err.Error())
			uf.Glog(&gibs, uf.GlogStruct{
				Level: uf.LevelWarn,
				Code:  "user.user.101",
				Msg:   "failed to create aim " + gibs.UfId.String(),
			})

			rb.AddError(AimCreateError)
			rb.BuildResponse(ctx, &out)
			return
		}
		aimInfo, err = LAimGetBySurferId(&gibs, in.SurferId)
		if err != nil {
			uf.Error(err.Error())
			uf.Glog(&gibs, uf.GlogStruct{
				Level: uf.LevelWarn,
				Code:  "user.user.102",
				Msg:   "failed to cleate aim " + gibs.UfId.String(),
			})

			rb.AddError(AimCreateError)
			rb.BuildResponse(ctx, &out)
			return
		}
	}

	chat, err := ut.LChatGetById(&gibs, in.ChatId)
	if err != nil {
		rb.AddError(GetChatError)
	}

	surfer, ures := uxt.SurferGetById(&gibs, in.SurferId)
	if ures.Errored() {
		rb.AddError(uxt.SurferNotFoundError)
	}

	msgInId, _ := gocql.RandomUUID()

	aimInfo.Iterating = true
	aimInfo.Surfer = surfer
	aimInfo.Cart = &CartPg{}
	aimInfo.Chat = chat
	aimInfo.ChatPlatform = in.ChatPlatform
	aimInfo.Recipient = in.Recipient
	aimInfo.MakoResponse = &MakoResponse{}
	aimInfo.MsgIn = ut.MsgPg{
		ChatId:         chat.ChatId,
		SenderSurferId: in.SurferId,
		ChatPlatform:   in.ChatPlatform,
		Content:        in.Content,
		MsgId:          msgInId,
	}
	aimInfo.MsgOut = ut.MsgPg{
		ChatId:         chat.ChatId,
		SenderSurferId: uf.FridayyUuid,
		Content:        "",
		ChatPlatform:   in.ChatPlatform,
	}
	aimInfo.QueryFilters = &QueryFilters{}
	aimInfo.Query = &QueryPg{}

	aimInfo.MsgIn.LChatSaveNewMsg(&gibs)

	ures, err = aimInfo.Process(&gibs)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "user.user.103",
			Interface: ures.Errors,
		})
	}

	rb.BuildResponse(ctx, &out)
}
