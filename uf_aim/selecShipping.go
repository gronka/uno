package uf_aim

import (
	"strconv"

	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/zinc"
)

type ShippingChoices struct {
	CheapestCost  int
	CheapestDays  int
	CheapestIndex int
	FastestCost   int
	FastestDays   int
	FastestIndex  int
	// Balanced choice is cheapest taking less than 6 days
	BalancedCost  int
	BalancedDays  int
	BalancedIndex int
}

func (sc *ShippingChoices) FindAndSetShippingChoices(porb zinc.ProductOffersResultBody) {
	index, choice := porb.GetCheapestChoice()
	sc.FastestIndex = index
	sc.FastestCost = choice.Price
	sc.FastestDays = choice.HandlingDays.Max

	index, choice = porb.GetFastestChoice()
	sc.FastestIndex = index
	sc.FastestCost = choice.Price
	sc.FastestDays = choice.HandlingDays.Max

	// the index will be set to -1 if no match is found
	index, choice = porb.GetDefaultChoice()
	sc.BalancedIndex = index
	sc.BalancedCost = choice.Price
	sc.BalancedDays = choice.HandlingDays.Max
}

func (sc *ShippingChoices) ShowShippingChoicesText() string {
	cheapestUsd := uf.IntToPriceString(sc.GetCheapestShippingPrice(99999))
	fastestUsd := uf.IntToPriceString(sc.GetFastestShippingPrice(99999))
	balancedUsd := uf.IntToPriceString(sc.GetBalancedShippingPrice(99999))

	//sc := cart.ShippingChoices
	cheapestText := "1) cheapest: " + cheapestUsd + ", will arrive in " +
		strconv.Itoa(sc.CheapestDays) + " days.\n"

	fastestText := "2) fastest: " + fastestUsd + ", will arrive in " +
		strconv.Itoa(sc.FastestDays) + " days.\n"

	balancedText := "3) balanced: " + balancedUsd + ", will arrive in " +
		strconv.Itoa(sc.BalancedDays) + " days.\n"

	return "Select a shipping and processing option:\n" + cheapestText +
		fastestText + balancedText
}

func (product *ProductPg) GetShippingChoices(gibs *uf.Gibs) (sc ShippingChoices) {
	if product.ScrapeEngine == "zinc" {
		productOffersResultBody, err := zinc.GetProductOffers(
			gibs.Conf,
			product.ScrapeEngineProductId,
			"amazon")
		if err != nil {
			uf.Error(err)
		}

		sc.FindAndSetShippingChoices(productOffersResultBody)
	} else {
		panic("invalid scrape engine")
	}

	return sc
}

func (sc *ShippingChoices) GetCheapestShippingPrice(itemPrice int) int {
	subtotal := itemPrice + sc.CheapestCost
	return int(float32(subtotal)*0.1) + 100
}

func (sc *ShippingChoices) GetFastestShippingPrice(itemPrice int) int {
	subtotal := itemPrice + sc.FastestCost
	return int(float32(subtotal)*0.1) + 100
}

func (sc *ShippingChoices) GetBalancedShippingPrice(itemPrice int) int {
	subtotal := itemPrice + sc.BalancedCost
	return int(float32(subtotal)*0.1) + 100
}
