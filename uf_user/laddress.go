package uf_user

import (
	"strings"

	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/zinc"
)

type AddressPg struct {
	AddressId    gocql.UUID
	SurferId     gocql.UUID
	IsBuilder    bool
	Name         string
	AddressLine1 string
	AddressLine2 string
	Postal       string
	Zip4         string
	City         string
	State        string
	Country      string
	Phone        string
	Instructions string
	IsValidated  bool
	TimeCreated  int64
	TimeUpdated  int64
}

func (ad *AddressPg) AsLetterhead() string {
	line2 := ad.AddressLine2
	if line2 != "" {
		line2 += "\n"
	}

	return ad.Name + "\n" +
		ad.AddressLine1 + "\n" +
		line2 +
		ad.City + ", " + ad.State + " " + ad.Postal
}

func (address *AddressPg) AsZincAddressObject() (obj zinc.AddressObject) {
	names := strings.SplitN(address.Name, " ", 2)
	lastName := ""
	if len(names) == 1 || names[1] == "" {
		lastName = "."
	}

	obj.FirstName = names[0]
	obj.LastName = lastName
	obj.AddressLine1 = address.AddressLine1
	obj.AddressLine2 = address.AddressLine2
	obj.ZipCode = address.Postal
	obj.City = address.City
	obj.State = address.State
	obj.Country = address.Country
	obj.PhoneNumber = address.Phone
	obj.Instructions = address.Instructions
	return
}

func (address *AddressPg) LAddressDelete(gibs *uf.Gibs) error {
	uf.Trace("delete address")
	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `DELETE FROM addresses 
	WHERE address_id = $1`)
	uf.FlashError(err)
	return err
}

func (address *AddressPg) LAddressUpsert(gibs *uf.Gibs) error {
	uf.Trace("insert address")
	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `INSERT INTO addresses (
		address_id,
		surfer_id, 
		is_builder,
		is_validated,
		name, 
		address_line_1,
		address_line_2,
		postal,
		zip4,
		city,
		state,
		country,
		phone,
		instructions,
		time_created,
		time_updated
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	ON CONFLICT (address_id) DO UPDATE SET
		is_builder = EXCLUDED.is_builder,
		is_validated = EXCLUDED.is_validated,
		name = EXCLUDED.name,
		address_line_1 = EXCLUDED.address_line_1,
		address_line_2 = EXCLUDED.address_line_2,
		postal = EXCLUDED.postal,
		zip4 = EXCLUDED.zip4,
		city = EXCLUDED.city,
		state = EXCLUDED.state,
		country = EXCLUDED.country,
		phone = EXCLUDED.phone,
		instructions = EXCLUDED.instructions,
		time_updated = EXCLUDED.time_updated`,

		address.AddressId,
		address.SurferId,
		address.IsBuilder,
		address.IsValidated,
		address.Name,
		address.AddressLine1,
		address.AddressLine2,
		address.Postal,
		address.Zip4,
		address.City,
		address.State,
		address.Country,
		address.Phone,
		address.Instructions,
		uf.NowStamp(),
		uf.NowStamp(),
	)
	uf.FlashError(err)

	return err
}

func LAddressGetById(gibs *uf.Gibs, addressId gocql.UUID) (address AddressPg, err error) {
	uf.Trace("select address")
	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
		address_id,
		surfer_id, 
		is_builder,
		is_validated,
		name, 
		address_line_1,
		address_line_2,
		postal,
		zip4,
		city,
		state,
		country,
		phone,
		instructions,
		time_created,
		time_updated
	FROM addresses WHERE address_id=$1
	ORDER BY time_created DESC
	LIMIT 1`,
		addressId,
	).Scan(
		&address.AddressId,
		&address.SurferId,
		&address.IsBuilder,
		&address.IsValidated,
		&address.Name,
		&address.AddressLine1,
		&address.AddressLine2,
		&address.Postal,
		&address.Zip4,
		&address.City,
		&address.State,
		&address.Country,
		&address.Phone,
		&address.Instructions,
		&address.TimeCreated,
		&address.TimeUpdated,
	)
	uf.FlashError(err)

	return
}

func LAddressGetNonBuilderBySurferId(gibs *uf.Gibs, surferId gocql.UUID) (address AddressPg, err error) {
	uf.Trace("select address")

	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
		address_id,
		surfer_id, 
		is_builder,
		is_validated,
		name, 
		address_line_1,
		address_line_2,
		postal,
		zip4,
		city,
		state,
		country,
		phone,
		instructions,
		time_created,
		time_updated
	FROM addresses WHERE surfer_id=$1 AND is_builder = $2
	ORDER BY time_created DESC
	LIMIT 1`,
		surferId,
		false,
	).Scan(
		&address.AddressId,
		&address.SurferId,
		&address.IsBuilder,
		&address.IsValidated,
		&address.Name,
		&address.AddressLine1,
		&address.AddressLine2,
		&address.Postal,
		&address.Zip4,
		&address.City,
		&address.State,
		&address.Country,
		&address.Phone,
		&address.Instructions,
		&address.TimeCreated,
		&address.TimeUpdated,
	)
	uf.FlashError(err)

	return
}

func LAddressGetBuilderForSurfer(gibs *uf.Gibs, surferId gocql.UUID) (address AddressPg, err error) {
	uf.Trace("select address")

	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
		address_id,
		surfer_id, 
		is_builder,
		is_validated,
		name, 
		address_line_1,
		address_line_2,
		postal,
		zip4,
		city,
		state,
		country,
		phone,
		instructions,
		time_created,
		time_updated
	FROM addresses WHERE surfer_id=$1 AND is_builder = $2
	ORDER BY time_created DESC
	LIMIT 1`,
		surferId,
		true,
	).Scan(
		&address.AddressId,
		&address.SurferId,
		&address.IsBuilder,
		&address.IsValidated,
		&address.Name,
		&address.AddressLine1,
		&address.AddressLine2,
		&address.Postal,
		&address.Zip4,
		&address.City,
		&address.State,
		&address.Country,
		&address.Phone,
		&address.Instructions,
		&address.TimeCreated,
		&address.TimeUpdated,
	)
	uf.FlashError(err)

	return
}
