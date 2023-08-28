package uf_public

import (
	"github.com/gocql/gocql"
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uxt"
)

type UfIdStruct struct {
	UfId gocql.UUID
}

type GlogCollection struct {
	Collection []uf.GlogStruct
}

func HGlogsGetByUfId(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HGlogGetByUfId called")
	in, out := UfIdStruct{}, GlogCollection{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	uf.Debug(gibs.ToddId)
	//TODO make admin only
	//rb.ExitIfPolicyFails(uf.PolicyToddIsSuperAdmin(gibs.ToddId))
	rb.ExitIfPolicyFails(uf.PolicyPublic())

	_ = uf.MakeRequest(
		&gibs,
		pile.Conf.AimAddress,
		uxt.PublicGlogsGetByUfIdV1,
		in,
		&out,
	)

	rb.BuildResponse(ctx, &out)
}
