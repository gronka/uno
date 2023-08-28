package uf_order

import (
	"net/http"

	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
	"gitlab.com/textfridayy/uno/zinc"
)

func ExecuteCreateOrder(
	conf *uf.Config,
	surfer uf_user.SurferPg,
	shippingAddress zinc.AddressObject,
	billingAddress *zinc.AddressObject,
	shippingInfo zinc.ShippingObject,
	products []zinc.ProductObject,
	maxPrice int,
	testing bool,
) (zinc.OrderCreateData, zinc.RequestIdBody, *http.Response, error) {
	//TODO idempotency key; we can probably use our internal OrderId
	order := zinc.OrderCreateData{}

	order.MaxPrice = uf.AddOurMargin(maxPrice)

	order.Products = products
	order.Retailer = zinc.RetailerAmazon
	order.ShippingAddress = shippingAddress

	order.Shipping = shippingInfo

	if billingAddress == nil {
		order.BillingAddress = getTylerBillingAddressObject()
	}

	order.PaymentMethod = zinc.PaymentMethodObject{
		NameOnCard:      conf.ZincCardName,
		Number:          conf.ZincCardNumber,
		ExpirationMonth: conf.ZincCardExpirationMonth,
		ExpirationYear:  conf.ZincCardExpirationYear,
		SecurityCode:    conf.ZincCardSecurityCode,
		UseGift:         false,
	}

	order.RetailerCredentials = zinc.RetailerCredentialsObject{
		Email:    conf.AmazonEmail,
		Password: conf.AmazonPassword,
	}

	FridayyOrderId, _ := gocql.RandomUUID()
	order.ClientNotes = zinc.ClientNotesMeta{
		FridayyOrderId: FridayyOrderId,
		SurferId:       surfer.SurferId,
		//QueryId             gocql.UUID
		//TODO: title
		Title: "first product name",
		//TODO:
		StripeCustomerId: surfer.StripeCustomerId,
		//TODO:
		StripePaymentMethod: surfer.StripeDefaultPaymentMethod,
	}

	order.Webhooks = zinc.WebhooksObject{
		RequestSucceeded: conf.OrderAddress + "/zinc.orderSuccess",
		RequestFailed:    conf.OrderAddress + "/zinc.orderFail",
	}

	orderResponse := zinc.RequestIdBody{}
	httpResponse, err := zinc.DoRequest(
		conf,
		zinc.PathOrders,
		uf.HttpMethodPost,
		order,
		map[string]string{},
		&orderResponse,
	)
	return order, orderResponse, httpResponse, err
}
