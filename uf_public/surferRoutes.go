package uf_public

import (
	"github.com/gronka/tg"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
	"gitlab.com/textfridayy/uno/uxt"
)

func HSurferGetById(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HSurferGetById called")
	in, out := SurferIdStruct{}, uf_user.SurferPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyToddIsSuperAdmin(gibs.ToddId))

	out, ures := uxt.SurferGetById(&gibs, in.SurferId)
	rb.AddErrors(ures.Errors)

	//TOOD: sanitizeo
	out.Sanitize()
	rb.BuildResponse(ctx, &out)
}
