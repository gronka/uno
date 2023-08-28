package uf_public

import (
	"github.com/gocql/gocql"
	"github.com/gronka/tg"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uxt"
)

type EmailPasswordStruct struct {
	Email    string
	Password string
}

type PhonePasswordStruct struct {
	Phone    string
	Password string
}

type SignInOut struct {
	NewJwt     string
	SessionJwt string //TODO
	//StripCustomerId     string //TODO
	SurferId gocql.UUID
}

func HSignInWithEmail(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HSurferSignInWithEmail called")
	in, out := EmailPasswordStruct{}, SignInOut{}
	gibs, rb := uf.InitGibsForWeb(ctx, &pile, &in)
	uf.Debug("checking policy")
	rb.ExitIfPolicyFails(uf.PolicyPublic())
	uf.Debug("policy passed")

	//log.Debug().Str("email", in.Email).Str("password", in.Password).Msg("HSurferSignInWithEmail")

	uf.Debug("get surfer by email")
	uf.Debug(in)
	surfer, ures := uxt.SurferGetByEmail(&gibs, in.Email)
	uf.Debug("surfer got")
	rb.AddErrors(ures.Errors)
	uf.Debug("added errors")
	if rb.HasErrors() {
		rb.BuildResponse(ctx, uf.EmptyStruct{})
		return
	}
	uf.Debug("still going")

	success := false
	if surfer.Password != "" && surfer.Password == in.Password {
		uf.Debug("signing in with")
		uf.Debug(surfer.SurferId)
		var err error
		out.NewJwt, err = uf.CreateJwt(&pile.Conf, surfer.SurferId, "email")
		if err != nil {
			uf.Error("failed to create jwt")
			panic(err)
		}
		out.SessionJwt = out.NewJwt
		out.SurferId = surfer.SurferId
		success = true
	} else {
		rb.AddError(CredentialsError)
	}

	log.Debug().Bool("isSignInSuccessful", success).Msg("hAccountSignInWithEmail")
	rb.BuildResponse(ctx, &out)
	return
}

func HSignInWithPhone(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	tg.Debug("HSurferSignInWithPhone called")
	in, out := PhonePasswordStruct{}, SignInOut{}
	gibs, rb := uf.InitGibsForWeb(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyPublic())

	surfer, ures := uxt.SurferGetByPhone(&gibs, in.Phone)
	rb.AddErrors(ures.Errors)
	if rb.HasErrors() {
		rb.BuildResponse(ctx, uf.EmptyStruct{})
		return
	}

	success := false
	if surfer.Password != "" && surfer.Password == in.Password {
		out.NewJwt, _ = uf.CreateJwt(&pile.Conf, surfer.SurferId, "phone")
		out.SessionJwt = out.NewJwt
		out.SurferId = surfer.SurferId
		success = true
	} else {
		rb.AddError(CredentialsError)
	}
	log.Debug().Bool("isSignInSuccessful", success).Msg("hAccountSignInWithEmail")

	rb.BuildResponse(ctx, &out)
	return
}
