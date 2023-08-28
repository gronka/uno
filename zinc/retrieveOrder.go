package zinc

import "gitlab.com/textfridayy/uno/uf"

type GetOrderOut struct {
	PriceComponents  PriceComponentsObject   `json:"price_components"`
	MerchantOrderIds []MerchantOrderIdObject `json:"merchant_order_ids"`
	Tracking         []TrackingObject        `json:"tracking"`
	Request          map[string]interface{}  `json:"request"`
	// there is no documentation for DeliveryDates on Zinc
	DeliveryDates []DeliveryDateObject `json:"delivery_dates"`
	// amazon only
	AccountStatus AccountStatusObject `json:"account_status"`

	// only in case the order is processing
	Type      string      `json:"_type"`
	Code      string      `json:"code"`
	Error     ErrorObject `json:"error"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

type PriceComponentsObject struct {
	Shipping              int                    `json:"shipping"`
	Products              []ProductInOrderObject `json:"products"`
	Subtotal              int                    `json:"subtotal"`
	Tax                   int                    `json:"tax"`
	Total                 int                    `json:"total"`
	GiftCertificate       int                    `json:"gift_certificate"`
	Currency              string                 `json:"currency"`
	PaymentCurrency       string                 `json:"payment_currency"`
	ConvertedPaymentTotal int                    `json:"converted_payment_total"`
}

// this object is for zinc, but not documented anywhere
type ProductInOrderObject struct {
	Price     int    `json:"price"`
	Quantity  int    `json:"quantity"`
	SellerId  string `json:"seller_id"`
	ProductId string `json:"product_id"`
}

type MerchantOrderIdObject struct {
	MerchantOrderId string   `json:"merchant_order_id"`
	Merchant        string   `json:"merchant"`
	Account         string   `json:"account"`
	PlacedAt        string   `json:"placed_at"`
	Tracking        []string `json:"tracking"`
	ProductIds      []string `json:"product_ids"`
	TrackingUrl     int      `json:"tracking_url"`
	DeliveryDate    int      `json:"delivery_date"`
}

type TrackingObject struct {
	MerchantOrderId string `json:"merchant_order_id"`
	// optional
	Carrier string `json:"carrier"`
	// optional
	TrackingNumber string `json:"tracking_number"`
	DeliveryStatus string `json:"delivery_status"`
	// optional
	TrackingUrl string   `json:"tracking_url"`
	ProductIds  []string `json:"product_ids"`
	// AmazonOnly
	RetailerTrackingNumber string `json:"retailer_tracking_number"`
	RetailerTrackingUrl    string `json:"retailer_tracking_url"`
	ObtainedAt             string `json:"obtained_at"`
}

// there is no documentation for this object
type DeliveryDateObject struct {
	DeliveryDate string `json:"delivery_date"`
	ProductId    string `json:"product_id"`
}

type AccountStatusObject struct {
	Prime    bool   `json:"prime"`
	Fresh    bool   `json:"fresh"`
	Business bool   `json:"business"`
	Charity  string `json:"charity"`
}

func GetOrderDetails(
	conf *uf.Config,
	orderRequestId string) (*GetOrderOut, error) {
	uf.Trace("zinc get order details for requestId: " + orderRequestId)

	out := GetOrderOut{}
	_, err := DoRequest(
		conf,
		string(PathOrders)+"/"+orderRequestId,
		uf.HttpMethodGet,
		nil,
		nil,
		&out,
	)
	uf.Trace("returning results")

	return &out, err
}
