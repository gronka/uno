package ut

import (
	"strconv"

	"github.com/gocql/gocql"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"gitlab.com/textfridayy/uno/uf"
)

type OrderPg struct {
	OrderId              gocql.UUID
	SurferId             gocql.UUID
	SurferSubscriptionId int
	BasketItemIds        []gocql.UUID

	PaymentAttemptIds []gocql.UUID
	WasPaymentMade    bool
	Status            string

	Title         string
	BasketPrice   int
	Tax           string
	Margin        int
	CreditToApply int
	Currency      string
	ShippingDays  int
	ShippingPrice int

	ShippingAddressId gocql.UUID
	ShippingName      string
	ShippingLine1     string
	ShippingLine2     string
	ShippingCity      string
	ShippingState     string
	ShippingPostal    string

	// Driver would be zinc or shopify
	Driver             string
	ZincOrderRequestId string

	TimePlaced  int64
	TimeCreated int64
	TimeUpdated int64
}

func (order *OrderPg) TotalPrice() int {
	tax, _ := strconv.ParseFloat(order.Tax, 10)
	taxedPrice := float64(order.BasketPrice) / 100.0 * tax
	taxedPriceInt := int(taxedPrice * 100)
	return taxedPriceInt + order.Margin
}

func (order *OrderPg) LUpsert(gibs *uf.Gibs) error {

	uf.Trace("upsert order")
	basketItemIdsArray := uf.UuidSliceToStringSlice(order.BasketItemIds)
	paymentAttemptIdsArray := uf.UuidSliceToStringSlice(order.PaymentAttemptIds)
	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `INSERT INTO orders (
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
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)
	ON CONFLICT (order_id) DO UPDATE SET
		basket_item_ids = EXCLUDED.basket_item_ids,

		payment_attempt_ids = EXCLUDED.payment_attempt_ids,
		was_payment_made = EXCLUDED.was_payment_made,
		status = EXCLUDED.status,

		title = EXCLUDED.title,
		basket_price = EXCLUDED.basket_price,
		tax = EXCLUDED.tax,
		margin = EXCLUDED.margin,
		credit_to_apply = EXCLUDED.credit_to_apply,
		currency = EXCLUDED.currency,
		shipping_days = EXCLUDED.shipping_days,
		shipping_price = EXCLUDED.shipping_price,

		shipping_address_id = EXCLUDED.shipping_address_id,
		shipping_name = EXCLUDED.shipping_name,
		shipping_line_1 = EXCLUDED.shipping_line_1,
		shipping_line_2 = EXCLUDED.shipping_line_2,
		shipping_city = EXCLUDED.shipping_city,
		shipping_state = EXCLUDED.shipping_state,
		shipping_postal = EXCLUDED.shipping_postal,

		driver = EXCLUDED.driver,
		zinc_order_request_id = EXCLUDED.zinc_order_request_id,

		time_placed = EXCLUDED.time_placed,
		time_updated = EXCLUDED.time_updated`,
		order.OrderId,
		order.SurferId,
		order.SurferSubscriptionId,
		pq.Array(basketItemIdsArray),

		pq.Array(paymentAttemptIdsArray),
		order.Status,

		order.Title,
		order.BasketPrice,
		order.Tax,
		order.Margin,
		order.CreditToApply,
		order.Currency,
		order.ShippingDays,
		order.ShippingPrice,

		order.ShippingAddressId,
		order.ShippingName,
		order.ShippingLine1,
		order.ShippingLine2,
		order.ShippingCity,
		order.ShippingState,
		order.ShippingPostal,

		order.Driver,
		order.ZincOrderRequestId,

		order.TimePlaced,
		order.TimeCreated,
		order.TimeUpdated,
	)
	uf.FlashError(err)

	return err
}

func LOrderSelectById(gibs *uf.Gibs, orderId gocql.UUID) (OrderPg, error) {

	order := OrderPg{}
	uf.Trace("select order")
	basketItemIdsArray := pgtype.UUIDArray{}
	paymentAttemptIdsArray := pgtype.UUIDArray{}
	err := gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
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
		orderId,
	).Scan(
		&order.OrderId,
		&order.SurferId,
		&order.SurferSubscriptionId,
		&basketItemIdsArray,

		&paymentAttemptIdsArray,
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
	uf.FlashError(err)

	order.BasketItemIds = uf.UuidArrayToUuidSlice(basketItemIdsArray)
	order.PaymentAttemptIds = uf.UuidArrayToUuidSlice(paymentAttemptIdsArray)

	return order, err
}

func LOrderSelectNewestBySurferId(gibs *uf.Gibs, surferId gocql.UUID) (OrderPg, error) {

	order := OrderPg{}
	uf.Trace("select order")
	basketItemIdsArray := pgtype.UUIDArray{}
	paymentAttemptIdsArray := pgtype.UUIDArray{}
	err := gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
		order_id,
		surfer_id,
		surfer_subscription_id,
		basket_item_ids,

		basket_item_ids,
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
	FROM orders WHERE surfer_id=$1 AND zinc_order_request_
	ORDER BY time_created DESC
	LIMIT 1`,
		surferId,
	).Scan(
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
	uf.FlashError(err)

	order.BasketItemIds = uf.UuidArrayToUuidSlice(basketItemIdsArray)
	order.PaymentAttemptIds = uf.UuidArrayToUuidSlice(paymentAttemptIdsArray)

	return order, err
}
