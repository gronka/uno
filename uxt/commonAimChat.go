package uxt

import (
	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
	"gitlab.com/textfridayy/uno/ut"
)

func ChatUpsert(
	gibs *uf.Gibs,
	chat *ut.ChatPg,
) uf.UfResponse {
	uf.Trace("ux.ChatUpsert")

	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.AimAddress,
		AimChatUpsertV1,
		chat,
		uf.EmptyStruct{},
	)

	return ures
}

func ChatGetBySurferId(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) (chat ut.ChatPg, ures uf.UfResponse) {
	uf.Trace("ux.ChatGetBySurferId")
	pkg := uf_user.SurferIdStruct{SurferId: surferId}

	ures = uf.MakeRequest(
		gibs,
		gibs.Conf.AimAddress,
		AimChatGetBySurferIdV1,
		pkg,
		&chat,
	)

	return
}

type ChatIdStruct struct {
	ChatId gocql.UUID
}

func ChatGetById(
	gibs *uf.Gibs,
	chatId gocql.UUID,
) (chat ut.ChatPg, ures uf.UfResponse) {
	uf.Trace("ux.ChatGetById")
	pkg := ChatIdStruct{ChatId: chatId}

	ures = uf.MakeRequest(
		gibs,
		gibs.Conf.AimAddress,
		AimChatGetByIdV1,
		pkg,
		&chat,
	)

	return
}
