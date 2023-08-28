package uf_order

import (
	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
)

type PaymentAttemptPg struct {
	PaymentAttemptId gocql.UUID
	OrderId          gocql.UUID

	// Driver is stripe or btcpay, etc
	Driver                string
	IsAutopaySelected     bool
	Status                string
	Notes                 string
	StripePaymentIntentId string
	StripePaymentMethodId string

	Price    int
	Currency string

	TermsAcceptedVersion string
	TimeCreated          int64
	TimeUpdated          int64
}

func (pa *PaymentAttemptPg) LUpsert(gibs *uf.Gibs) (err error) {

	uf.Trace("upsert payment_attempt")
	_, err = gibs.Pile.Pool.Exec(gibs.Ctx, `INSERT INTO payment_attempts (
		payment_attempt_id,
		order_id,

		driver,
		is_autopay_selected,
		status,
		notes,
		stripe_payment_intent_id,
		stripe_payment_method_id,

		price,
		currency,

		terms_accepted_version,
		time_created,
		time_updated
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	ON CONFLICT (payment_attempt_id) DO UPDATE SET
		order_id = EXCLUDED.order_id,

		driver = EXCLUDED.driver,
		is_autopay_selected = EXCLUDED.is_autopay_selected,
		status = EXCLUDED.status,
		notes = EXCLUDED.notes,
		stripe_payment_intent_id = EXCLUDED.stripe_payment_intent_id,
		stripe_payment_method_id = EXCLUDED.stripe_payment_method_id,

		price = EXCLUDED.price,
		currency = EXCLUDED.currency,

		terms_accepted_version = EXCLUDED.terms_accepted_version,
		time_updated = EXCLUDED.time_updated`,
		pa.PaymentAttemptId,
		pa.OrderId,

		pa.Driver,
		pa.IsAutopaySelected,
		pa.Status,
		pa.Notes,
		pa.StripePaymentIntentId,
		pa.StripePaymentMethodId,

		pa.Price,
		pa.Currency,

		pa.TermsAcceptedVersion,
		pa.TimeCreated,
		pa.TimeUpdated,
	)
	uf.FlashError(err)

	return err
}

func (pa *PaymentAttemptPg) LSelectById(gibs *uf.Gibs, paId gocql.UUID) error {

	uf.Trace("select payment_attempt")
	err := gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT (
		payment_attempt_id,
		order_id,

		driver,
		is_autopay_selected,
		status,
		notes,
		stripe_payment_intent_id,
		stripe_payment_method_id,

		price,
		currency,

		terms_accepted_version,
		time_created,
		time_updated
	) FROM payment_attempts WHERE payment_attempt_id=$1`,
		paId,
	).Scan(
		&pa.PaymentAttemptId,
		&pa.OrderId,

		&pa.Driver,
		&pa.IsAutopaySelected,
		&pa.Status,
		&pa.Notes,
		&pa.StripePaymentIntentId,
		&pa.StripePaymentMethodId,

		&pa.Price,
		&pa.Currency,

		&pa.TermsAcceptedVersion,
		&pa.TimeCreated,
		&pa.TimeUpdated,
	)
	uf.FlashError(err)

	return err
}

func LSelectPaymentAttemptNeedingService(gibs *uf.Gibs) (pa PaymentAttemptPg, err error) {

	uf.Trace("select payment_attempt")
	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT (
		payment_attempt_id,
		order_id,

		driver,
		is_autopay_selected,
		status,
		notes,
		stripe_payment_intent_id,
		stripe_payment_method_id,

		price,
		currency,

		terms_accepted_version,
		time_created,
		time_updated
	) FROM payment_attempts 
	WHERE status=$1 
	LIMIT 1`,
		"created",
	).Scan(
		&pa.PaymentAttemptId,
		&pa.OrderId,

		&pa.Driver,
		&pa.IsAutopaySelected,
		&pa.Status,
		&pa.Notes,
		&pa.StripePaymentIntentId,
		&pa.StripePaymentMethodId,

		&pa.Price,
		&pa.Currency,

		&pa.TermsAcceptedVersion,
		&pa.TimeCreated,
		&pa.TimeUpdated,
	)
	uf.FlashError(err)

	return pa, err
}
