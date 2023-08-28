package uxt

import (
	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
)

func SurferGetById(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) (uf_user.SurferPg, uf.UfResponse) {
	pkg := uf_user.SurferIdStruct{SurferId: surferId}

	surfer := uf_user.SurferPg{}
	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserSurferGetByIdV1,
		pkg,
		&surfer,
	)

	return surfer, ures
}

func SurferGetByEmail(
	gibs *uf.Gibs,
	email string,
) (uf_user.SurferPg, uf.UfResponse) {
	pkg := uf_user.EmailStruct{Email: email}

	surfer := uf_user.SurferPg{}
	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserSurferGetByEmailV1,
		pkg,
		&surfer,
	)

	return surfer, ures
}

func SurferGetByPhone(
	gibs *uf.Gibs,
	phone string,
) (uf_user.SurferPg, uf.UfResponse) {
	pkg := uf_user.PhoneStruct{Phone: phone}

	surfer := uf_user.SurferPg{}
	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserSurferGetByPhoneV1,
		pkg,
		&surfer,
	)

	return surfer, ures
}

func SurferCreateStripeCustomerId(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) uf.UfResponse {
	pkg := uf_user.SurferIdStruct{SurferId: surferId}

	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserSurferCreateStripeCustomerIdV1,
		pkg,
		uf.EmptyStruct{},
	)

	return ures
}

func SurferCreateFromPhoneOrEmail(
	gibs *uf.Gibs,
	name string,
	phoneOrEmail string,
) (uf_user.SurferPg, uf.UfResponse) {
	pkg := uf_user.PhoneOrEmailStruct{Name: name, PhoneOrEmail: phoneOrEmail}

	surfer := uf_user.SurferPg{}
	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserSurferGetOrCreateFromPhoneOrEmailV1,
		pkg,
		&surfer,
	)

	return surfer, ures
}
