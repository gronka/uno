package uf_public

import (
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_aim"
	"gitlab.com/textfridayy/uno/uf_user"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

func SendMsgToAim(
	chatPlatform ut.ChatPlatform,
	chat *ut.ChatPg,
	in *LoopHookIn,
	gibs *uf.Gibs,
	surfer *uf_user.SurferPg,
	rb *uf.ResponseBuilder,
	ures *uf.UfResponse) {

	rb.AddErrors(ures.Errors)

	//TODO: do something smarter than jump if ures.Errored()
	if uf.UuidIsZero(chat.ChatId) || ures.Errored() {
		chat = &ut.ChatPg{
			ChatId:         uf.RandomUUID(),
			SurferId:       surfer.SurferId,
			Teaser:         in.Text,
			TeaserSurferId: surfer.SurferId,
		}

		uxt.ChatUpsert(gibs, chat)
	}

	uf.Trace("4")
	pkg := uf_aim.MsgFromSurfer{
		SurferId:     surfer.SurferId,
		SurferPhone:  surfer.Phone,
		SurferEmail:  surfer.Email,
		Recipient:    in.Recipient,
		ChatId:       chat.ChatId,
		ChatPlatform: chatPlatform,
		Content:      in.Text,
		GroupId:      "",
	}

	uf.Trace("5")
	uf.Trace(gibs)
	uf.Trace(gibs.Pile.Conf.AimAddress)
	uf.Trace("5")
	uf.Trace(pkg)
	ures2 := uf.MakeRequest(
		gibs,
		gibs.Pile.Conf.AimAddress,
		uxt.AimMessageReceiveV1,
		pkg,
		uf.EmptyStruct{},
	)

	uf.Glog(gibs, uf.GlogStruct{
		Level:     uf.LevelDebug,
		ChatId:    chat.ChatId,
		SurferId:  surfer.SurferId,
		Code:      "public.101",
		Msg:       "Received message from user",
		Interface: chat,
	})

	rb.AddErrors(ures2.Errors)

	ures = &ures2
	uf.Trace("6")

}
