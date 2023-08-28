package uf_aim

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/lib/pq"
	"gitlab.com/textfridayy/uno/uf"
)

type ProductPg struct {
	ProductId gocql.UUID
	// QueryHash murmur3(noun, sorted(adjectives))
	QueryHash gocql.UUID
	// BadMatch allows us to prevent products from appearing
	BadMatch bool

	// taxonomy
	Noun                 string
	Category             string
	AdjectivesSorted     []string
	AntiAdjectivesSorted []string

	// details
	Title         string
	Description   string
	Price         int
	Currency      string
	ReviewCount   int
	Stars         float32
	ShippingDays  int
	ShippingPrice int

	// meta
	ScrapeEngine          string // zinc, etc
	ScrapeEngineProductId string // productId on the scrape engine
	StoreName             string
	StoreEngine           string
	StoreUrl              string
	ProductUrl            string
	ImageUrl              string

	TimeCreated int64
	TimeUpdated int64
}

func (product *ProductPg) GetMargin() int {
	// add $1 + 10% for Fridayy profit marging
	margin := int(100 + float32(product.Price)*0.1)
	return margin
}

func (product *ProductPg) LineItemPrice() int {
	margin := product.GetMargin()
	return margin + product.Price + product.ShippingPrice
}

func (product *ProductPg) PriceAsUsdString() string {
	return uf.IntToPriceString(product.LineItemPrice())
}

func (product *ProductPg) LProductUpsert(gibs *uf.Gibs) error {
	uf.Trace("insert product")
	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `INSERT INTO products (
		product_id,
		query_hash,
		bad_match,

		noun,
		category,
		adjectives_sorted,
		anti_adjectives_sorted,

		title,
		description,
		price,
		currency,
		review_count,
		stars,
		shipping_days,
		shipping_price,

		scrape_engine,
		scrape_engine_product_id,
		store_name,
		store_engine,
		store_url,
		product_url,
		image_url,

		time_created,
		time_updated
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24)
	ON CONFLICT (product_id) DO UPDATE SET
		query_hash = EXCLUDED.query_hash,
		bad_match = EXCLUDED.bad_match,

		noun = EXCLUDED.noun,
		category = EXCLUDED.category,
		adjectives_sorted = EXCLUDED.adjectives_sorted,
		anti_adjectives_sorted = EXCLUDED.anti_adjectives_sorted,

		title = EXCLUDED.title,
		description = EXCLUDED.description,
		price = EXCLUDED.price,
		currency = EXCLUDED.currency,
		review_count = EXCLUDED.review_count,
		stars = EXCLUDED.stars,
		shipping_days = EXCLUDED.shipping_days,
		shipping_price = EXCLUDED.shipping_price,

		scrape_engine = EXCLUDED.scrape_engine,
		scrape_engine_product_id = EXCLUDED.scrape_engine_product_id,
		store_name = EXCLUDED.store_name,
		store_engine = EXCLUDED.store_engine,
		store_url = EXCLUDED.store_url,
		product_url = EXCLUDED.product_url,
		image_url = EXCLUDED.image_url,

		time_updated = EXCLUDED.time_updated`,

		product.ProductId,
		product.QueryHash,
		product.BadMatch,

		product.Noun,
		product.Category,
		pq.Array(product.AdjectivesSorted),
		pq.Array(product.AntiAdjectivesSorted),
		//adjSortedPq,
		//antiAdjSortedPq,

		product.Title,
		product.Description,
		product.Price,
		product.Currency,
		product.ReviewCount,
		product.Stars,
		product.ShippingDays,
		product.ShippingPrice,

		product.ScrapeEngine,
		product.ScrapeEngineProductId,
		product.StoreName,
		product.StoreEngine,
		product.StoreUrl,
		product.ProductUrl,
		product.ImageUrl,

		uf.NowStamp(),
		uf.NowStamp(),
	)
	uf.FlashError(err)

	return err
}

func (query *QueryPg) FindOnUno(gibs *uf.Gibs) (product ProductPg, err error) {
	queryHash := query.ComputeQueryHash()

	minPrice := query.MinPrice
	if minPrice < query.AbsoluteMinPrice {
		minPrice = query.AbsoluteMinPrice
	}

	maxPrice := query.MaxPrice
	if maxPrice < query.AbsoluteMaxPrice {
		maxPrice = query.AbsoluteMaxPrice
	}

	minPriceMod := ""
	if minPrice != 0 {
		minPriceMod = fmt.Sprintf(" AND price > %d ", minPrice)
	}

	maxPriceMod := ""
	if maxPrice != 0 {
		maxPriceMod = fmt.Sprintf(" AND price < %d ", maxPrice)
	}

	uf.Trace("find product from uno")
	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT 
		product_id,
		query_hash,
		bad_match,

		noun,
		category,
		adjectives_sorted,
		anti_adjectives_sorted,

		title,
		description,
		price,
		currency,
		review_count,
		stars,
		shipping_days,
		shipping_price,

		scrape_engine,
		scrape_engine_product_id,
		store_name,
		store_engine,
		store_url,
		product_url,
		image_url,

		time_created,
		time_updated
	FROM products 
	WHERE query_hash=$1 AND bad_match=false `+
		minPriceMod+maxPriceMod+
		"LIMIT 1",
		queryHash,
	).Scan(
		&product.ProductId,
		&product.QueryHash,
		&product.BadMatch,

		&product.Noun,
		&product.Category,
		pq.Array(&product.AdjectivesSorted),
		pq.Array(&product.AntiAdjectivesSorted),

		&product.Title,
		&product.Description,
		&product.Price,
		&product.Currency,
		&product.ReviewCount,
		&product.Stars,
		&product.ShippingDays,
		&product.ShippingPrice,

		&product.ScrapeEngine,
		&product.ScrapeEngineProductId,
		&product.StoreName,
		&product.StoreEngine,
		&product.StoreUrl,
		&product.ProductUrl,
		&product.ImageUrl,

		&product.TimeCreated,
		&product.TimeUpdated,
	)
	uf.FlashError(err)

	uf.Warn(product)

	return
}

func LProductGetById(gibs *uf.Gibs, productId gocql.UUID) (
	product ProductPg,
	err error,
) {
	uf.Trace("select product " + productId.String())
	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT 
		product_id,
		query_hash,
		bad_match,

		noun,
		category,
		adjectives_sorted,
		anti_adjectives_sorted,

		title,
		description,
		price,
		currency,
		review_count,
		stars,
		shipping_days,
		shipping_price,

		scrape_engine,
		scrape_engine_product_id,
		store_name,
		store_engine,
		store_url,
		product_url,
		image_url,

		time_created,
		time_updated
	FROM products 
	WHERE product_id=$1`,
		productId,
	).Scan(
		&product.ProductId,
		&product.QueryHash,
		&product.BadMatch,

		&product.Noun,
		&product.Category,
		pq.Array(&product.AdjectivesSorted),
		pq.Array(&product.AntiAdjectivesSorted),

		&product.Title,
		&product.Description,
		&product.Price,
		&product.Currency,
		&product.ReviewCount,
		&product.Stars,
		&product.ShippingDays,
		&product.ShippingPrice,

		&product.ScrapeEngine,
		&product.ScrapeEngineProductId,
		&product.StoreName,
		&product.StoreEngine,
		&product.StoreUrl,
		&product.ProductUrl,
		&product.ImageUrl,

		&product.TimeCreated,
		&product.TimeUpdated,
	)
	uf.FlashError(err)

	return
}
