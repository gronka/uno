package uf_aim

import (
	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
)

type CartPg struct {
	CartId   gocql.UUID
	SurferId gocql.UUID

	ProductIds        []gocql.UUID
	Counts            []int64
	Prices            []int64
	ShippingAddressId gocql.UUID
	IsBuilderAddress  bool

	// purchasing members
	ShippingChoices *ShippingChoices
	ShippingChoice  string
	BuilderAddress  uf_user.AddressPg

	TimeCreated int
	TimeUpdated int
}

type CartItem struct {
	ProductId gocql.UUID
	Quantity  int
}

func (cart *CartPg) GetTotalProductPrices() (total int64) {
	for _, val := range cart.Prices {
		total += val
	}
	return
}

func (cart *CartPg) GetFinalPrice(gibs *uf.Gibs) (total int) {
	for _, productId := range cart.ProductIds {
		product, _ := LProductGetById(gibs, productId)
		total += product.LineItemPrice()
	}
	return
}

func (cart *CartPg) AddToCart(productId gocql.UUID, count, price int64) {
	if cart.ProductIds != nil {
		cart.ProductIds = append(cart.ProductIds, productId)
	} else {
		cart.ProductIds = []gocql.UUID{productId}
	}

	if cart.Counts != nil {
		cart.Counts = append(cart.Counts, count)
	} else {
		cart.Counts = []int64{count}
	}

	if cart.Prices != nil {
		cart.Prices = append(cart.Prices, price)
	} else {
		cart.Prices = []int64{price}
	}
}
