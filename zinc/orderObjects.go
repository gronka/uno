package zinc

import (
	"github.com/gocql/gocql"
)

type OrderCreateData struct {
	// Zinc variables
	Retailer        ZincRetailer    `json:"retailer,omitempty"`
	Products        []ProductObject `json:"products,omitempty"`
	ShippingAddress AddressObject   `json:"shipping_address,omitempty"`
	// shipping option. choices: cheapest, fastest, amazon_day, free (fails for
	// items without free shipping). Do not provide this if you Provide the
	// `shipping` attribute
	//Shipping string          `json:"shipping"`
	Shipping            ShippingObject            `json:"shipping,omitempty"`
	BillingAddress      AddressObject             `json:"billing_address,omitempty"`
	PaymentMethod       PaymentMethodObject       `json:"payment_method,omitempty"`
	RetailerCredentials RetailerCredentialsObject `json:"retailer_credentials,omitempty"`
	IsGift              bool                      `json:"is_gift,omitempty"`
	MaxPrice            int                       `json:"max_price,omitempty"`

	//
	// Optional parameters for all orders
	//
	GiftMessage string          `json:"gift_message,omitempty"`
	RequireGift bool            `json:"require_gift,omitempty"`
	Webhooks    WebhooksObject  `json:"webhooks,omitempty"`
	ClientNotes ClientNotesMeta `json:"client_notes,omitempty"`

	PromoCodes []PromoCodeObject `json:"promo_codes,omitempty"`
	// defaults to false
	StrictExpiredPoductId bool `json:"strict_expired_poduct_id,omitempty"`
	// Amazon business accounts only
	PoNumber int `json:"po_number,omitempty"`
	// Amazon only. If ship_method == amazon_day, here you specify an exact name.
	AmazonDay string `json:"amazon_day,omitempty"`
	// defaults to false
	FailIfTaxed           bool   `json:"fail_if_taxed,omitempty"`
	MaxDeliveryDays       int    `json:"max_delivery_days,omitempty"`
	TakeBuyboxOffers      bool   `json:"take_buybox_offers,omitempty"`
	ForceOffersPostalCode string `json:"force_offers_postal_code,omitempty"`
}

type PromoCodeObject struct {
	Code string `json:"code,omitempty"`
	// optional
	Optional bool `json:"optional,omitempty"`
	// optional
	MerchantId string `json:"Merchantmd,_imitempty"`
	// optional; default 0; requires merchant_id to be submitted
	DiscountAmount int `json:"discount_amount,omitempty"`
	// optional; default 0; requires merchant_id to be submitted. Percentage
	// between 0 and 100.
	DiscountPercentage int `json:"discount_percentage,omitempty"`
	// optional; requires merchant_id to be submitted
	CostOverride int `json:"cost_override,omitempty"`
}

type WebhooksObject struct {
	RequestSucceeded string `json:"request_succeeded,omitempty"`
	RequestFailed    string `json:"request_failed,omitempty"`
	TrackingObtained string `json:"tracking_obtained,omitempty"`
	StatusUpdated    string `json:"status_updated,omitempty"`
	// ZMA only
	//CaseUpdated    string `json:"case_updated"`
}

type PaymentMethodObject struct {
	NameOnCard                string `json:"name_on_card,omitempty"`
	Number                    string `json:"number,omitempty"`
	SecurityCode              string `json:"security_code,omitempty"`
	ExpirationMonth           int    `json:"expiration_month,omitempty"`
	ExpirationYear            int    `json:"expiration_year,omitempty"`
	UseGift                   bool   `json:"use_gift"`
	UseAccountPaymentDefaults bool   `json:"use_account_payment_defaults,omitempty"`
	IsVirtualCard             bool   `json:"is_virtual_card,omitempty"`
}
type RetailerCredentialsObject struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	// optional
	VerificationCode string `json:"verification_code,omitempty"`
	// optional
	Totp2faKey string `json:"totp_2fa_key,omitempty"`
}

type ShippingObject struct {
	// ordering of shipping methods. Options are `price` or `speed`
	OrderBy  string `json:"order_by,omitempty"`
	MaxDays  int    `json:"max_days,omitempty"`
	MaxPrice int    `json:"max_price,omitempty"`
}

type AddressObject struct {
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	ZipCode      string `json:"zip_code,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Country      string `json:"country,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

type ProductObject struct {
	ProductId               string                        `json:"product_id,omitempty"`
	Quantity                int                           `json:"quantity,omitempty"`
	SellerSelectionCriteria SellerSelectionCriteriaObject `json:"seller_selection_criteria,omitempty"`
}

type SellerSelectionCriteriaObject struct {
	// list of allowed conditions. Example: ["New"]
	ConditionIn []string `json:"condition_in,omitempty"`
	// list of disabled conditions. Example: ["New"]
	ConditionNotIn   []string `json:"condition_not_in,omitempty"`
	FirstPartySeller bool     `json:"first_party_seller,omitempty"`
	HandlingDaysMax  int      `json:"handling_days_max,omitempty"`
	International    bool     `json:"international,omitempty"`
	MaxItemPrice     int      `json:"max_item_price,omitempty"`

	// amazon only
	Addon bool `json:"addon,omitempty"`
	// amazon only
	BuyBox bool `json:"buy_box,omitempty"`
	// amazon only
	MerchantIdIn []string `json:"merchant_id_in,omitempty"`
	// amazon only
	MerchantIdNotIn []string `json:"merchant_id_not_in,omitempty"`
	// amazon only
	MinSellerNumRatings bool `json:"min_seller_num_ratings,omitempty"`
	// amazon only
	MinSellerPercentPositiveFeedback int `json:"min_seller_percent_positive_feedback,omitempty"`
	// amazon only
	Prime bool `json:"prime,omitempty"`
	// amazon only
	AllowOos bool `json:"allow_oos,omitempty"`
}

// ClientNotesMeta can contain whatever data we like
type ClientNotesMeta struct {
	FridayyOrderId      gocql.UUID
	SurferId            gocql.UUID
	QueryId             gocql.UUID
	Title               string
	StripeCustomerId    string
	StripePaymentMethod string
}
