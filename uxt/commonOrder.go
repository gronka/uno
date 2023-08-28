package uxt

import (
	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
	"gitlab.com/textfridayy/uno/ut"
)

func OrderCreate(
	gibs *uf.Gibs,
	orderIn ut.OrderCreateIn,
) (orderOut ut.OrderCreateOut, ures uf.UfResponse) {

	ures = uf.MakeRequest(
		gibs,
		gibs.Conf.OrderAddress,
		OrderOrderCreateV1,
		orderIn,
		&orderOut,
	)

	return orderOut, ures
}

func PaymentAttemptCreate(
	gibs *uf.Gibs,
	pa ut.PaymentAttemptCreateIn,
) (paOut ut.PaymentAttemptCreateOut, ures uf.UfResponse) {

	ures = uf.MakeRequest(
		gibs,
		gibs.Conf.OrderAddress,
		OrderPaymentAttemptCreateV1,
		pa,
		&paOut,
	)

	return paOut, ures
}

func OrderGetBySurferIdNewest(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) (order ut.OrderPg, ures uf.UfResponse) {
	pkg := uf_user.SurferIdStruct{SurferId: surferId}

	ures = uf.MakeRequest(
		gibs,
		gibs.Conf.OrderAddress,
		OrderOrderCreateV1,
		pkg,
		&order,
	)

	return order, ures
}

func OrderGetDetails(
	gibs *uf.Gibs,
	orderId gocql.UUID,
) (details ut.OrderDetails, ures uf.UfResponse) {
	pkg := uf_user.SurferIdStruct{SurferId: orderId}

	ures = uf.MakeRequest(
		gibs,
		gibs.Conf.OrderAddress,
		OrderOrderCreateV1,
		pkg,
		&details,
	)

	return details, ures
}
