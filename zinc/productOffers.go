package zinc

import (
	"math"

	"gitlab.com/textfridayy/uno/uf"
)

type ProductOffersResultBody struct {
	// uno variables
	CheapestOfferPos int
	DefaultOfferPos  int
	FastestOfferPos  int

	// Status can be processing (if async: true was set), "failed" or "completed"
	Status    string               `json:"status"`
	Retailer  string               `json:"retailer"`
	Offers    []ProductOfferResult `json:"offers"`
	TimeStamp int                  `json:"timestamp"`
}

// Details at https://docs.zincapi.com/#product-offers
// WARNING: poorly documented by Zinc
type ProductOfferResult struct {
	// true if product is an addon and must be purchased as a bundle
	Addon bool `json:"addon"`
	// Condition can be: New, Refurbished, Used - Like New, Used - Very Good,
	// Used - Good, Used - Acceptable, Unnacceptable
	Condition    string             `json:"condition"`
	HandlingDays HandlingDaysObject `json:"handling_days"`

	International bool `json:"international"`
	// amazon only: unique ID for an item sold by any merchant on Amazon
	OfferId              string                  `json:"offer_id"`
	Price                int                     `json:"price"`
	MarketplaceFulfilled bool                    `json:"marketplace_fulfilled"`
	Seller               SellerObject            `json:"seller"`
	ShippingOptions      []ShippingOptionsObject `json:"shipping_options"`

	MaxAge int `json:"status"`
	// amazon only
	PrimeOnly bool `json:"prime_only"`
	// costco only
	MemberOnly bool `json:"member_only"`
	// ND amazon only, example: "B00KFP6NHO"
	Asin string `json:"asin"`

	// ND
	Currency string `json:"currency"`
	// ND fulfillment by amazon?
	FbaBadge bool `json:"fba_badge"`
	// ND ?
	PrimeBadge bool `json:"prime_badge"`
	// ND shipping description
	Greytext string `json:"greytext"`
	// ND
	KnownGreytext bool `json:"known_greytext"`
	// ND lol nice
	Available bool `json:"available"`
	// ND
	EstimatedOfferCount int `json:"estimated_offer_count"`
	// ND
	MinimumQuantity int `json:"minimum_quantity"`
	// ND
	ExpiredProductId bool `json:"expired_product_id"`

	//---- more ND
	//Comments ??? `json:"comments"`
	//CouponDiscountFixed ??? `json:"coupon_discount_fixed"`
	//CouponDiscountPercent ??? `json:"coupon_discount_percent"`
	//BuyBoxWinner ??? `json:"buy_box_winner"`
}

type HandlingDaysObject struct {
	Max int `json:"max"`
	Min int `json:"min"`
}

type DeliveryDaysObject struct {
	Max int `json:"max"`
	Min int `json:"min"`
}

type SellerObject struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	PercentPositive int    `json:"percent_positive"`
	NumRatings      int    `json:"num_ratings"`
	FirstParty      bool   `json:"first_party"`
}

type ShippingOptionsObject struct {
	Name         string             `json:"name"`
	Price        int                `json:"price"`
	DeliveryDays DeliveryDaysObject `json:"delivery_days"`
	Greytext     string             `json:"greytext"`
}

func (porb *ProductOffersResultBody) GetFastestChoice() (int, ProductOfferResult) {
	smallestValue := math.MaxInt32
	for index, offer := range porb.Offers {
		if offer.HandlingDays.Max < smallestValue {
			porb.FastestOfferPos = index
		}
	}

	if smallestValue == math.MaxInt32 {
		porb.FastestOfferPos = -1
		return porb.FastestOfferPos, ProductOfferResult{}
	} else {
		return porb.FastestOfferPos, porb.Offers[porb.FastestOfferPos]
	}
}

func (porb *ProductOffersResultBody) GetCheapestChoice() (int, ProductOfferResult) {
	smallestValue := math.MaxInt32
	for index, offer := range porb.Offers {
		if offer.Price < smallestValue {
			porb.CheapestOfferPos = index
		}
	}

	if smallestValue == math.MaxInt32 {
		porb.CheapestOfferPos = -1
		return porb.CheapestOfferPos, ProductOfferResult{}
	} else {
		return porb.CheapestOfferPos, porb.Offers[porb.CheapestOfferPos]
	}
}

func (porb *ProductOffersResultBody) GetDefaultChoice() (int, ProductOfferResult) {
	smallestValue := math.MaxInt32
	for index, offer := range porb.Offers {
		if offer.HandlingDays.Max <= 6 {
			if offer.Price < smallestValue {
				porb.DefaultOfferPos = index
			}
		}
	}

	if smallestValue == math.MaxInt32 {
		porb.DefaultOfferPos = -1
		return porb.DefaultOfferPos, ProductOfferResult{}
	} else {
		return porb.DefaultOfferPos, porb.Offers[porb.DefaultOfferPos]
	}
}

func GetProductOffers(
	conf *uf.Config,
	productId string,
	retailer string) (ProductOffersResultBody, error) {

	path := PathProducts + "/" + productId + "/offers"
	out := ProductOffersResultBody{}
	_, err := DoRequest(
		conf,
		path,
		uf.HttpMethodGet,
		nil,
		map[string]string{
			"retailer": retailer,
		},
		&out,
	)

	return out, err
}
