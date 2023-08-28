package uf_order

import (
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/checkout/session"
	"github.com/stripe/stripe-go/v73/customer"
	"github.com/stripe/stripe-go/v73/paymentintent"
	"github.com/stripe/stripe-go/v73/paymentmethod"
	"gitlab.com/textfridayy/uno/uf"
)

func GetPaymentMethods(gibs *uf.Gibs, stripeCustomerId string) {
	stripe.Key = gibs.Conf.StripeSecretKey
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(stripeCustomerId),
		Type:     stripe.String("card"),
	}
	i := paymentmethod.List(params)
	for i.Next() {
		pm := i.PaymentMethod()
		fmt.Println(pm)
		fmt.Println(pm.Card.Last4)
	}
}

func GetStripeCustomer(gibs *uf.Gibs, stripeCustomerId string) *stripe.Customer {
	stripe.Key = gibs.Conf.StripeSecretKey
	cust, err := customer.Get(stripeCustomerId, nil)
	uf.FlashError(err)
	return cust
}

func GetDefaultPaymentMethod(gibs *uf.Gibs, stripeCustomerId string) *stripe.PaymentMethod {
	cust := GetStripeCustomer(gibs, stripeCustomerId)
	method := &stripe.PaymentMethod{}
	var err error
	if cust.DefaultSource != nil && cust.DefaultSource.ID != "" {
		method, err = paymentmethod.Get(cust.DefaultSource.ID, nil)
		uf.FlashError(err)
	}
	return method
}

func IsPaymentMethodExpired(pm *stripe.PaymentMethod) bool {
	if pm.Card == nil {
		return false
	}

	year, month, _ := time.Now().Date()
	if int64(year) > pm.Card.ExpYear {
		return true
	} else if int64(year) == pm.Card.ExpYear {
		if int64(month) > pm.Card.ExpMonth {
			return true
		}
	}

	return false
}

func CreateStripeSession(
	gibs *uf.Gibs,
	stripeCustomerId string,
	price int,
	title string,
) (*stripe.CheckoutSession, error) {
	stripe.Key = gibs.Conf.StripeSecretKey

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(stripeCustomerId),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Mode:     stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(title),
					},
					UnitAmount: stripe.Int64(int64(price)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			SetupFutureUsage: stripe.String("off_session"),
		},
		SuccessURL: stripe.String(gibs.Conf.AppDn + "/payment.success"),
		CancelURL:  stripe.String(gibs.Conf.AppDn + "/payment.cancel"),
	}

	sesh, err := session.New(params)
	uf.FlashError(err)

	return sesh, err
}

func CreateStripePIWithSavedCard(
	gibs *uf.Gibs,
	stripeCustomerId string,
	price int,
	queryString string,
	pmId string,
) (*stripe.PaymentIntent, error) {
	stripe.Key = gibs.Conf.StripeSecretKey

	params := &stripe.PaymentIntentParams{
		Customer:         stripe.String(stripeCustomerId),
		SetupFutureUsage: stripe.String("off_session"),
		Amount:           stripe.Int64(int64(price)),
		Currency:         stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		StatementDescriptor: stripe.String(queryString),
		PaymentMethod:       stripe.String(pmId),
	}

	pi, err := paymentintent.New(params)
	uf.FlashError(err)

	return pi, err
}

func GetCheckoutSession(
	gibs *uf.Gibs,
	sessionId string,
) (*stripe.CheckoutSession, error) {
	stripe.Key = gibs.Conf.StripeSecretKey

	sesh, err := session.Get(sessionId, nil)
	uf.FlashError(err)

	return sesh, err
}
