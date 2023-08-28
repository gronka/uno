package uf_aim

import (
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/zinc"
)

func (aimInfo *AimInfo) PerformSearchInitial(gibs *uf.Gibs) (product ProductPg) {
	uf.Trace("PerformSearch")
	aimInfo.CreateNewQuery()

	// search local database
	product, err := aimInfo.Query.FindOnUno(gibs)
	if err != nil {
		uf.Error(err)
	}

	uf.Trace("was product found?")
	if product.Title == "" {
		uf.Trace("no")
		// if no result, use zinc
		searchResultBody, err := zinc.SearchForProduct(
			gibs.Conf,
			aimInfo.MakoResponse.AsQueryString(),
			"amazon")
		if err != nil {
			uf.Error(err)
		} else {
			for _, sr := range searchResultBody.Results {
				// add all new results to our database
				newResult := ZincSearchResultToProduct(sr, aimInfo.Query)
				err = newResult.LProductUpsert(gibs)
				if err != nil {
					uf.Error(err)
				}
			}

			// perform search on uno again
			product, err = aimInfo.Query.FindOnUno(gibs)
			if err != nil {
				uf.Error(err)
			}
		}
	}

	aimInfo.Query.CurrentProductId = product.ProductId
	aimInfo.Query.CurrentProductPrice = product.Price

	return
}

func (aimInfo *AimInfo) PerformSearchRefine(gibs *uf.Gibs) (product ProductPg) {
	// search local database
	product, err := aimInfo.Query.FindOnUno(gibs)
	if err != nil {
		uf.Error(err)
	}

	if product.Title == "" {
		product, err := LProductGetById(gibs, aimInfo.Query.CurrentProductId)
		if err != nil {
			uf.Error(err)
		}
		aimInfo.ShowProductResponse(gibs, product, "Search failed. Still viewing:\n")
	} else {
		aimInfo.Query.CurrentProductId = product.ProductId
		aimInfo.Query.CurrentProductPrice = product.Price
		aimInfo.ShowProductResponse(gibs, product, "")
	}
	return
}

func (aimInfo *AimInfo) MakeBranchSearchInitiateResponse(gibs *uf.Gibs, product ProductPg) {
	out := &aimInfo.MsgOut
	if aimInfo.Error != nil {
		out.Code = "search_item_failed"
		out.Content = "We're sorry, but your search has failed. Please try again."
		return
	}

	if product.Title == "" {
		out.Code = "search_item_no_results"
		out.Content = "Sorry, I could not find anything. Try again?"
		return
	}

	aimInfo.ShowProductResponse(gibs, product, "")
}

func (aimInfo *AimInfo) ShowProductResponse(gibs *uf.Gibs, product ProductPg, prefix string) {
	if product.ShippingDays == 0 {
		sc := product.GetShippingChoices(gibs)
		product.ShippingDays = sc.BalancedDays
		product.ShippingPrice = sc.BalancedCost
		product.LProductUpsert(gibs)
	}

	out := &aimInfo.MsgOut
	out.Code = "search_item_result"
	out.MediaUrl = product.ImageUrl
	out.Content = "Name: " + product.Title + "\n" +
		"Price: " + product.PriceAsUsdString() + "\n" +
		//"Adjectives: " + strings.Join(product.AdjectivesSorted, ", ") + "\n" +
		"What do you think?\n" +
		"1) purchase this item\n" +
		"2) for a cheaper choice\n" +
		"3) for a higher quality choice\n" +
		"9) cancel this search"
}

func (aimInfo *AimInfo) MakeBranchSearchPurchaseResponse() {
	out := &aimInfo.MsgOut
	out.Code = "search_item_purchase"
	out.Content = "Let's try and buy!"
}
