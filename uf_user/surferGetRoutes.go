package uf_user

import (
	"strings"

	"github.com/gocql/gocql"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
)

type PhoneStruct struct {
	Phone string
}

type PhoneOrEmailStruct struct {
	Name         string
	PhoneOrEmail string
}

type EmailStruct struct {
	Email string
}

type SurferIdStruct struct {
	SurferId gocql.UUID
}

func HSurferGetOrCreateFromPhoneOrEmail(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HSurferGetOrCreateFromPhoneOrEmail")
	in, out := PhoneOrEmailStruct{}, SurferPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Debug(in.PhoneOrEmail)
	isPhone := false
	if strings.Contains(in.PhoneOrEmail, "@") {
		isPhone = false
	} else {
		isPhone = true
	}

	var err error
	if isPhone {
		err = pile.Pool.QueryRow(ctx,
			"SELECT surfer_id, email, phone, password FROM surfers WHERE phone=$1",
			in.PhoneOrEmail,
		).Scan(
			&out.SurferId,
			&out.Email,
			&out.Phone,
			&out.Password,
		)
	} else {
		err = pile.Pool.QueryRow(ctx,
			"SELECT surfer_id, email, phone, password FROM surfers WHERE phone=$1",
			in.PhoneOrEmail,
		).Scan(
			&out.SurferId,
			&out.Email,
			&out.Phone,
			&out.Password,
		)
	}

	if err != nil {
		uf.Debug("got the surfer: " + out.SurferId.String())
		uf.Debug(err.Error())
		if err.Error() == "no rows in result set" {
			uf.Debug("no user for phoneOrEmail: " + in.PhoneOrEmail)

			if isPhone {
				err = lCreateUserFromPhone(ctx, pile, in.PhoneOrEmail)
				out, _ = LSurferGetByPhone(&gibs, in.PhoneOrEmail)
			} else {
				err = lCreateUserFromEmail(ctx, pile, in.PhoneOrEmail)
				out, _ = LSurferGetByEmail(&gibs, in.PhoneOrEmail)
			}

			if err != nil {
				rb.AddError(CreateUserError)
			}
		} else {
			uf.Debug("select failed for phone: " + in.PhoneOrEmail)
			rb.AddError(GetUserError)
		}
	}

	rb.BuildResponse(ctx, out)
}

func HSurferGetBySurferId(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HSurferGetBySurferId coming in")
	in, out := SurferIdStruct{}, SurferPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	out, err := LSurferGetById(&gibs, in.SurferId)
	if err != nil {
		uf.Debug("select failed for SurferId: " + in.SurferId.String())
		uf.Debug(err.Error())
		rb.AddError(GetUserError)
	}

	rb.BuildResponse(ctx, out)
}

func HSurferGetByEmail(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HSurferGetByEmail coming in")
	in, out := EmailStruct{}, SurferPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Debug(in.Email)

	out, err := LSurferGetByEmail(&gibs, in.Email)
	if err != nil {
		uf.Debug("select failed for Email: " + in.Email)
		uf.Debug(err.Error())
		rb.AddError(GetUserError)
	}

	rb.BuildResponse(ctx, out)
}

func HSurferGetByPhone(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HSurferGetByPhone coming in")
	in, out := PhoneStruct{}, SurferPg{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Debug(in.Phone)

	out, err := LSurferGetByPhone(&gibs, in.Phone)
	if err != nil {
		uf.Debug("select failed for Phone: " + in.Phone)
		uf.Debug(err.Error())
		rb.AddError(GetUserError)
	}

	rb.BuildResponse(ctx, out)
}
