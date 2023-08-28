package zinc

type ZincPath string

const basePath = "https://api.zinc.io/v1"

const (
	PathCancellations ZincPath = basePath + "/cancellations"
	PathSearch                 = basePath + "/search"
	PathProducts               = basePath + "/products"
	PathOrders                 = basePath + "/orders"
)

type Epid struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type EpidsMap map[string]string

type PackageDimensions struct {
	Weight WeightDimension `json:"weight"`
	Size   Size            `json:"size"`
}

type WeightDimension struct {
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Size struct {
	Width  SizeDimension `json:"width"`
	Depth  SizeDimension `json:"depth"`
	Length SizeDimension `json:"length"`
}

type SizeDimension struct {
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type VariantSpecific struct {
	// example: size
	Dimension float64 `json:"dimension"`
	// example: "1", "2"
	Value string `json:"value"`
}

type VariantSpecificsPerProduct struct {
	VariantSpecifics []VariantSpecific `json:"variant_specifics"`
	ProductId        string            `json:"product_id"`
}

type ZincRetailer string

const (
	RetailerAmazon = "amazon"
)

type RequestIdBody struct {
	RequestId string `json:"request_id"`
}
