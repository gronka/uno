package zinc

import (
	"gitlab.com/textfridayy/uno/uf"
)

// Details at https://docs.zincapi.com/?shell#product-details
type ProductResult struct {
	Status           string                       `json:"status"`
	Retailer         string                       `json:"retailer"`
	ProductId        string                       `json:"product_id"`
	Timestamp        int                          `json:"timestamp"`
	Title            string                       `json:"title"`
	ProductDetails   []string                     `json:"product_details"`
	FeatureBullets   []string                     `json:"feature_bullets"`
	Brand            string                       `json:"brand"`
	MainImage        string                       `json:"main_image"`
	Images           []string                     `json:"images"`
	VariantSpecifics []VariantSpecific            `json:"variant_specifics"`
	AllVariants      []VariantSpecificsPerProduct `json:"all_variants"`
	Categories       []string                     `json:"categories"`

	ProductDescription string            `json:"product_description"`
	Epids              []Epid            `json:"epids"`
	EpidsMap           EpidsMap          `json:"epids_map"`
	PackageDimensions  PackageDimensions `json:"package_dimensions"`

	// only for books
	Authors []string `json:"authors"`

	// aliexpress only
	ItemLocation string `json:"item_location"`
	// amazon only
	OriginalRetailPrice int `json:"original_retail_price"`
	// amazon only
	Price int `json:"price"`
	// amazon only
	ReviewCount int `json:"review_count"`
	// amazon only
	Stars float64 `json:"stars"`
	// amazon only
	QuestionCount int `json:"question_count"`
	// amazon only, example: "B00KFP6NHO"
	Asin string `json:"asin"`
	// amazon only
	Fresh bool `json:"fresh"`
	// amazon only
	Pantry bool `json:"pantry"`
	// amazon only
	Handmade bool `json:"handmade"`
	// amazon only
	Digital bool `json:"digital"`
	// amazon only
	BuyapiHint bool `json:"buyapi_hint"`
	// CostCo only
	ItemNumber string `json:"item_number"`
}

func GetProductDetails(
	conf *uf.Config,
	productId string,
	retailer string) (ProductResult, error) {

	path := PathProducts + "/" + productId
	out := ProductResult{}
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
