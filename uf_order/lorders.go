package uf_order

import (
	"github.com/gocql/gocql"
	"github.com/jackc/pgtype"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/ut"
)

type OrderPgCollection struct {
	Collection []ut.OrderPg
}

func LOrdersSelectBySurferId(gibs *uf.Gibs, surferId gocql.UUID) (OrderPgCollection, error) {

	uf.Trace("select orders by surferId")
	rows, err := gibs.Pile.Pool.Query(gibs.Ctx, `SELECT
		order_id,
		surfer_id,
		surfer_subscription_id,
		basket_item_ids,

		payment_attempt_ids,
		was_payment_made,
		status,

		title,
		basket_price,
		tax,
		margin,
		credit_to_apply,
		currency,
		shipping_days,
		shipping_price,

		shipping_address_id,
		shipping_name,
		shipping_line_1,
		shipping_line_2,
		shipping_city,
		shipping_state,
		shipping_postal,

		driver,
		zinc_order_request_id,

		time_placed,
		time_created,
		time_updated
	 FROM orders WHERE order_id=$1`,
		surferId,
	)
	defer rows.Close()

	out := OrderPgCollection{}
	out.Collection = make([]ut.OrderPg, 0)

	for rows.Next() {
		basketItemIdsArray := pgtype.UUIDArray{}
		paymentAttemptIdsArray := pgtype.UUIDArray{}

		order := ut.OrderPg{}
		err = rows.Scan(
			&order.OrderId,
			&order.SurferId,
			&order.SurferSubscriptionId,
			&basketItemIdsArray,

			&paymentAttemptIdsArray,
			&order.WasPaymentMade,
			&order.Status,

			&order.Title,
			&order.BasketPrice,
			&order.Tax,
			&order.Margin,
			&order.CreditToApply,
			&order.Currency,
			&order.ShippingDays,
			&order.ShippingPrice,

			&order.ShippingAddressId,
			&order.ShippingName,
			&order.ShippingLine1,
			&order.ShippingLine2,
			&order.ShippingCity,
			&order.ShippingState,
			&order.ShippingPostal,

			&order.Driver,
			&order.ZincOrderRequestId,

			&order.TimePlaced,
			&order.TimeCreated,
			&order.TimeUpdated,
		)
		order.BasketItemIds = uf.UuidArrayToUuidSlice(basketItemIdsArray)
		order.PaymentAttemptIds = uf.UuidArrayToUuidSlice(paymentAttemptIdsArray)

		out.Collection = append(out.Collection, order)
	}
	uf.FlashError(err)

	return out, err
}
