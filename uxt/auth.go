package uxt

import (
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
)

type PhoneStruct struct {
	Phone string
}

func AuthViaLoopHook(
	gibs *uf.Gibs,
	pile uf.Pile,
	isPhone bool,
	phoneOrEmail string,
) (surfer uf_user.SurferPg, ures uf.UfResponse) {

	uf.Trace("phoneOrEmail: " + phoneOrEmail)
	if isPhone {
		surfer, ures = SurferGetByPhone(gibs, phoneOrEmail)
	} else {
		surfer, ures = SurferGetByEmail(gibs, phoneOrEmail)
	}

	authToken, err := uf.CreateJwt(&pile.Conf, surfer.SurferId, "loop")
	if err != nil {
		ures.AddError(FailedToCreateJwtError)
	}
	gibs.UfAuth = authToken

	return surfer, ures
}

func AuthViaSipHook(
	gibs *uf.Gibs,
	pile uf.Pile,
	phone string,
) (surfer uf_user.SurferPg, ures uf.UfResponse) {

	surfer, ures = SurferGetByPhone(gibs, phone)

	authToken, err := uf.CreateJwt(&pile.Conf, surfer.SurferId, "sip")
	if err != nil {
		ures.AddError(FailedToCreateJwtError)
	}
	gibs.UfAuth = authToken

	return surfer, ures
}
