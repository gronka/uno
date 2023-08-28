package uf_aim

import (
	"github.com/gocql/gocql"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"gitlab.com/textfridayy/uno/uf"
)

func (cart *CartPg) LCartUpsert(gibs *uf.Gibs) error {
	//if cart.ShippingChoices == nil {
	//cart.ShippingChoices = &ShippingChoices{
	//CheapestCost: -1,
	//CheapestDays: -1,
	//FastestCost:  -1,
	//FastestDays:  -1,
	//BalancedCost: -1,
	//BalancedDays: -1,
	//}
	//}

	uf.Trace("upsert cart")
	productIdsArray := uf.UuidSliceToStringSlice(cart.ProductIds)
	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `INSERT INTO carts (
		cart_id,
		surfer_id,
		product_ids,
		counts,
		shipping_address_id,
		is_builder_address,

		time_created,
		time_updated
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (cart_id) DO UPDATE SET
		product_ids = EXCLUDED.product_ids,
		counts = EXCLUDED.counts,
		shipping_address_id = EXCLUDED.shipping_address_id,
		is_builder_address = EXCLUDED.is_builder_address,

		time_updated = EXCLUDED.time_updated`,

		cart.CartId,
		cart.SurferId,
		pq.Array(productIdsArray),
		pq.Array(cart.Counts),
		cart.ShippingAddressId,
		cart.IsBuilderAddress,

		uf.NowStamp(),
		uf.NowStamp(),
	)
	uf.FlashError(err)

	return err
}

func (cart *CartPg) LCartGetCurrent(gibs *uf.Gibs, surferId gocql.UUID) error {

	uf.Trace("select cart " + cart.CartId.String())
	productIdsArray := pgtype.UUIDArray{}
	err := gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
		cart_id,
		surfer_id,
		product_ids,
		counts,
		shipping_address_id,
		is_builder_address,

		time_created,
		time_updated
	FROM carts WHERE surfer_id=$1
	ORDER BY time_created DESC
	LIMIT 1`,
		surferId,
	).Scan(
		&cart.CartId,
		&cart.SurferId,
		&productIdsArray,
		(*pq.Int64Array)(&cart.Counts),
		&cart.ShippingAddressId,
		&cart.IsBuilderAddress,

		&cart.TimeCreated,
		&cart.TimeUpdated,
	)
	uf.FlashError(err)

	cart.ProductIds = uf.UuidArrayToUuidSlice(productIdsArray)

	return err
}
