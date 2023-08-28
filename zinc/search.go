package zinc

import (
	"fmt"

	"gitlab.com/textfridayy/uno/uf"
)

type ZincRequest struct {
	Path string
}

// details at https://docs.zincapi.com/?shell#product-search
type SearchResultBody struct {
	// Status can be processing (if async: true was set), "failed" or "completed"
	Status   string         `json:"status"`
	Retailer string         `json:"retailer"`
	Results  []SearchResult `json:"results"`
}

// details at https://docs.zincapi.com/?shell#product-search
type SearchResult struct {
	ProductId  string  `json:"product_id"`
	Title      string  `json:"title"`
	Image      string  `json:"image"`
	NumReviews int     `json:"num_reviews"`
	Stars      float32 `json:"stars"`
	Fresh      bool    `json:"fresh"`
	Price      int     `json:"price"`
}

func (sr *SearchResult) PriceAsString() string {
	return uf.IntToPriceString(sr.Price)
}

func (sr *SearchResult) PPrint() string {
	return fmt.Sprintf(
		"ProductId: %v\nTitle: %v\nImage: %v\nNumReviews: %v\n"+
			"Stars: %f\nFresh: %v\nPrice: %v\n",
		sr.ProductId,
		sr.Title,
		sr.Image,
		sr.NumReviews,
		sr.Stars,
		sr.Fresh,
		sr.Price)
}

func (body *SearchResultBody) PPrintFirst() string {
	if len(body.Results) > 0 {
		return body.Results[0].PPrint()
	}
	return "No results"
}

func SearchForProduct(
	conf *uf.Config,
	query,
	retailer string) (*SearchResultBody, error) {
	uf.Trace("zinc searching for query: " + query)

	out := SearchResultBody{}
	_, err := DoRequest(
		conf,
		string(PathSearch),
		uf.HttpMethodGet,
		nil,
		map[string]string{
			"page":     "1",
			"query":    query,
			"retailer": retailer,
		},
		&out,
	)
	uf.Trace("returning results")

	return &out, err
}
