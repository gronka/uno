package ut

import (
	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
)

type BasketItemPg struct {
	BasketItemId gocql.UUID
	OrderId      gocql.UUID
	ProductId    gocql.UUID
	Title        string
	Quantity     int
	UnitPrice    int
	Currency     string

	ScrapeEngine          string // zinc, etc
	ScrapeEngineProductId string // productId on the scrape engine
	StoreUrl              string
	ProductUrl            string
	ImageUrl              string

	TimeCreated int64
	TimeUpdated int64
}

func (item *BasketItemPg) LUpsert(gibs *uf.Gibs) error {

	uf.Trace("upsert basket_item")
	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `INSERT INTO basket_items (
		basket_item_id,
		order_id,
		product_id,
		title,
		quantity,
		unit_price,
		currency,

		scrape_engine,
		scrape_engine_product_id,
		store_url,
		product_url,
		image_url,

		time_created,
		time_updated
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	ON CONFLICT (basket_item_id) DO UPDATE SET
		order_id = EXCLUDED.order_id,
		product_id = EXCLUDED.product_id,
		title = EXCLUDED.title,
		quantity = EXCLUDED.quantity,
		unit_price = EXCLUDED.unit_price,
		currency = EXCLUDED.currency,

		scrape_engine = EXCLUDED.scrape_engine,
		scrape_engine_product_id = EXCLUDED.scrape_engine_product_id,
		store_url = EXCLUDED.store_url,
		product_url = EXCLUDED.product_url,
		image_url = EXCLUDED.image_url,
		time_updated = EXCLUDED.time_updated`,

		item.BasketItemId,
		item.OrderId,
		item.ProductId,
		item.Title,
		item.Quantity,
		item.UnitPrice,
		item.Currency,

		item.ScrapeEngine,
		item.ScrapeEngineProductId,
		item.StoreUrl,
		item.ProductUrl,
		item.ImageUrl,

		uf.NowStamp(),
		uf.NowStamp(),
	)
	uf.FlashError(err)

	return err
}

func (item *BasketItemPg) LSelectById(gibs *uf.Gibs) error {

	uf.Trace("select basket_item")
	err := gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
		basket_item_id,
		order_id,
		product_id,
		title,
		quantity,
		unit_price,
		currency,

		scrape_engine,
		scrape_engine_product_id,
		store_url,
		product_url,
		image_url,

		time_created,
		time_updated
	FROM basket_items WHERE basket_item_id=$1
	LIMIT 1`,
		item.BasketItemId,
	).Scan(
		&item.BasketItemId,
		&item.OrderId,
		&item.ProductId,
		&item.Title,
		&item.Quantity,
		&item.UnitPrice,
		&item.Currency,

		&item.ScrapeEngine,
		&item.ScrapeEngineProductId,
		&item.StoreUrl,
		&item.ProductUrl,
		&item.ImageUrl,

		item.TimeCreated,
		item.TimeUpdated,
	)
	uf.FlashError(err)

	return err
}
