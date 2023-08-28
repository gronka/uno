package uf_aim

import (
	"github.com/gocql/gocql"
	"github.com/lib/pq"
	"gitlab.com/textfridayy/uno/uf"
)

func LQueryGetById(gibs *uf.Gibs, queryId gocql.UUID) (query QueryPg, err error) {
	arr := make([]string, 0)
	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT 
		query_id,
		surfer_id,
		current_product_id,
		skipped_product_ids,
		current_product_price,
		absolute_max_price,
		absolute_min_price,
		count,
		
		noun,
		category,
		adjectives_sorted,
		anti_adjectives_sorted,

		time_created,
		time_updated
	FROM queries WHERE query_id=$1`,
		queryId,
	).Scan(
		&query.QueryId,
		&query.SurferId,
		&query.CurrentProductId,
		pq.Array(&arr),
		&query.CurrentProductPrice,
		&query.AbsoluteMaxPrice,
		&query.AbsoluteMinPrice,
		&query.Count,

		&query.Noun,
		&query.Category,
		pq.Array(&query.AdjectivesSorted),
		pq.Array(&query.AntiAdjectivesSorted),

		&query.TimeCreated,
		&query.TimeUpdated,
	)

	query.SkippedProductIds = uf.StringSliceToUuidSlice(arr)
	uf.FlashError(err)

	return
}

func (query *QueryPg) LQueryUpsert(gibs *uf.Gibs) error {
	arr := uf.UuidSliceToStringSlice(query.SkippedProductIds)

	uf.Trace("insert query")
	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `INSERT INTO queries (
		query_id,
		surfer_id,
		current_product_id,
		skipped_product_ids,
		current_product_price,
		absolute_max_price,
		absolute_min_price,
		count,

		noun,
		category,
		adjectives_sorted,
		anti_adjectives_sorted,
		time_created,
		time_updated
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	ON CONFLICT (query_id) DO UPDATE SET
		current_product_id = EXCLUDED.current_product_id,
		skipped_product_ids = EXCLUDED.skipped_product_ids,
		current_product_price = EXCLUDED.current_product_price,
		absolute_max_price = EXCLUDED.absolute_max_price,
		absolute_min_price = EXCLUDED.absolute_min_price,
		count = EXCLUDED.count,

		noun = EXCLUDED.noun,
		category = EXCLUDED.category,
		adjectives_sorted = EXCLUDED.adjectives_sorted,
		anti_adjectives_sorted = EXCLUDED.anti_adjectives_sorted,
		time_updated = EXCLUDED.time_updated`,

		query.QueryId,
		query.SurferId,
		query.CurrentProductId,
		pq.Array(arr),
		query.CurrentProductPrice,
		query.AbsoluteMaxPrice,
		query.AbsoluteMinPrice,
		query.Count,

		query.Noun,
		query.Category,
		pq.Array(query.AdjectivesSorted),
		pq.Array(query.AntiAdjectivesSorted),
		uf.NowStamp(),
		uf.NowStamp(),
	)
	uf.FlashError(err)

	return err
}
