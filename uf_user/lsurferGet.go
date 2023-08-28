package uf_user

import (
	"context"

	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
)

type SurferPg struct {
	SurferId                   gocql.UUID
	Email                      string
	Phone                      string
	Password                   string
	Name                       string
	IsEmailVerified            bool
	IsPhoneVerified            bool
	StripeCustomerId           string
	StripeDefaultPaymentMethod string
	//DefaultAdddressId gocql.UUID // might not be needed
	TimeCreated int64
	TimeUpdated int64
}

func (su *SurferPg) String() string {
	return "SurferId: " + su.SurferId.String() +
		", Email: " + su.Email +
		", Name: " + su.Name +
		", Phone: " + su.Phone
}

func (surfer *SurferPg) Sanitize() error {
	surfer.Password = ""
	return nil
}

func lCreateUserFromPhone(
	ctx context.Context,
	pile uf.Pile,
	phone string) error {
	uf.Trace("CreateUserFromPhone")

	_, err := pile.Pool.Exec(ctx, `INSERT INTO surfers (
	surfer_id,
	email,
	phone) VALUES (
	$1,
	'',
	$2)`,
		uf.RandomUUID(),
		phone)
	uf.FlashError(err)

	return err
}

func lCreateUserFromEmail(
	ctx context.Context,
	pile uf.Pile,
	email string) error {
	uf.Trace("CreateUserFromEmail")

	_, err := pile.Pool.Exec(ctx, `INSERT INTO surfers (
	surfer_id,
	email,
	phone) VALUES (
	$1,
	$2,
	'')`,
		uf.RandomUUID(),
		email)
	uf.FlashError(err)

	return err
}

func LSurferGetById(gibs *uf.Gibs, surferId gocql.UUID) (surfer SurferPg, err error) {
	uf.Trace("select surfer by id: " + surferId.String())

	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
		surfer_id, 
		email,
		phone,
		password, 
		name, 
		is_email_verified,
		is_phone_verified,
		stripe_customer_id, 
		time_created,
		time_updated
	FROM surfers WHERE surfer_id=$1`,
		surferId,
	).Scan(
		&surfer.SurferId,
		&surfer.Email,
		&surfer.Phone,
		&surfer.Password,
		&surfer.Name,
		&surfer.IsEmailVerified,
		&surfer.IsPhoneVerified,
		&surfer.StripeCustomerId,
		&surfer.TimeCreated,
		&surfer.TimeUpdated,
	)
	uf.FlashError(err)

	return
}

func LSurferGetByPhone(gibs *uf.Gibs, phone string) (surfer SurferPg, err error) {
	uf.Trace("select surfer by phone: " + phone)

	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
		surfer_id, 
		email,
		phone,
		password, 
		name, 
		is_email_verified,
		is_phone_verified,
		stripe_customer_id, 
		time_created,
		time_updated
	FROM surfers WHERE phone=$1`,
		phone,
	).Scan(
		&surfer.SurferId,
		&surfer.Email,
		&surfer.Phone,
		&surfer.Password,
		&surfer.Name,
		&surfer.IsEmailVerified,
		&surfer.IsPhoneVerified,
		&surfer.StripeCustomerId,
		&surfer.TimeCreated,
		&surfer.TimeUpdated,
	)
	uf.FlashError(err)

	return
}

func LSurferGetByEmail(gibs *uf.Gibs, email string) (surfer SurferPg, err error) {
	uf.Trace("select surfer by email: " + email)

	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
		surfer_id, 
		email,
		phone,
		password, 
		name, 
		is_email_verified,
		is_phone_verified,
		stripe_customer_id, 
		time_created,
		time_updated
	FROM surfers WHERE email=$1`,
		email,
	).Scan(
		&surfer.SurferId,
		&surfer.Email,
		&surfer.Phone,
		&surfer.Password,
		&surfer.Name,
		&surfer.IsEmailVerified,
		&surfer.IsPhoneVerified,
		&surfer.StripeCustomerId,
		&surfer.TimeCreated,
		&surfer.TimeUpdated,
	)
	uf.FlashError(err)

	return
}
