package uxt

import (
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
)

func SendReply(
	gibs *uf.Gibs,
	pkg *ut.MsgPg,
) uf.UfResponse {
	uf.Trace("send reply")

	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.BorderAddress,
		BorderSendMsgOutV1,
		pkg,
		uf.EmptyStruct{},
	)

	return ures
}
