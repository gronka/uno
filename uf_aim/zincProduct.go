package uf_aim

import (
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/zinc"
)

func ZincSearchResultToProduct(
	sr zinc.SearchResult,
	query *QueryPg,
) (product ProductPg) {
	product.ProductId = uf.RandomUUID()
	product.QueryHash = query.ComputeQueryHash()
	product.Noun = query.Noun
	product.Category = query.Category
	product.AdjectivesSorted = query.AdjectivesSorted
	product.AntiAdjectivesSorted = query.AntiAdjectivesSorted

	product.Title = sr.Title
	//product.Description = sr.Description
	product.ImageUrl = sr.Image
	product.Price = sr.Price
	product.Currency = "USD"
	product.ReviewCount = sr.NumReviews
	product.Stars = sr.Stars

	product.ScrapeEngine = "zinc"
	product.ScrapeEngineProductId = sr.ProductId
	//product.StoreName = sr.StoreName
	product.StoreEngine = "amazon"
	//product.StoreName = sr.StoreName
	//product.StoreUrl = sr.StoreUrl
	//product.ProductUrl = sr.ProductUrl
	product.ImageUrl = sr.Image

	return product
}
