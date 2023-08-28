package ut

import "github.com/gocql/gocql"

type OrderCreateIn struct {
	SurferId          gocql.UUID
	Title             string
	Price             int
	Margin            int
	Currency          string
	ShippingDays      int
	ShippingPrice     int
	ShippingAddressId gocql.UUID
	BasketItems       []BasketItemPg
}

type OrderCreateOut struct {
	OrderId                gocql.UUID
	StripeCustomerId       string
	PaymentMethodId        string
	IsPaymentMethodExpired bool
	CardLast4              string
}

type PaymentAttemptCreateIn struct {
	OrderId      gocql.UUID
	SurferId     gocql.UUID
	UseSavedCard bool
}

type PaymentAttemptCreateOut struct {
	PaymentUrl       string
	PaymentSucceeded bool
}
