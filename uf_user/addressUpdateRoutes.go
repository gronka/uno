package uf_user

import (
	"github.com/gocql/gocql"
	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
)

type AddressUpdateNameStruct struct {
	AddressId gocql.UUID
	Name      string
}

func HAddressUpdateName(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressUpdateName coming in")
	in, out := AddressUpdateNameStruct{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Debug(in)

	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `UPDATE addresses SET
		name = $1
		WHERE address_id = $2`,
		in.Name,
		in.AddressId,
	)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "user.address.101",
			Interface: err,
		})
	}

	rb.BuildResponse(ctx, out)
}

type AddressUpdatePostalStruct struct {
	AddressId gocql.UUID
	Postal    string
}

var CityStateLookupFailedError = uf.ApiError{
	Code: "city_state_lookup_failed_error",
	Msg:  "Failed to find City/State from Zip",
}

func HAddressUpdatePostalPlus(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressUpdatePostalPlus coming in")
	in, out := AddressUpdatePostalStruct{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Debug(in)

	lookup := UspsZipcodeToCityState(&gibs, in.Postal)
	if lookup.ZipC.Zip5 != in.Postal {
		rb.AddError(CityStateLookupFailedError)
	} else {
		_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `UPDATE addresses SET
		postal = $1,
		city = $2,
		state = $3
		WHERE address_id = $4`,
			in.Postal,
			lookup.ZipC.City,
			lookup.ZipC.State,
			in.AddressId,
		)
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "user.address.102",
			Interface: err,
		})
	}

	rb.BuildResponse(ctx, out)
}

type AddressUpdateAddressLine1Struct struct {
	AddressId    gocql.UUID
	AddressLine1 string
}

func HAddressUpdateAddressLine1(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressUpdateAddressLine1 coming in")
	in, out := AddressUpdateAddressLine1Struct{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Debug(in)

	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `UPDATE addresses SET
		address_line_1 = $1
		WHERE address_id = $2`,
		in.AddressLine1,
		in.AddressId,
	)
	uf.Glog(&gibs, uf.GlogStruct{
		Level:     uf.LevelError,
		Code:      "user.address.103",
		Interface: err,
	})

	rb.BuildResponse(ctx, out)
}

type AddressUpdateAddressLine2Struct struct {
	AddressId    gocql.UUID
	AddressLine2 string
}

func HAddressUpdateAddressLine2(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressUpdateAddressLine2 coming in")
	in, out := AddressUpdateAddressLine2Struct{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	uf.Debug(in)

	line := in.AddressLine2
	if line == "none" || line == "None" {
		line = ""
	}

	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `UPDATE addresses SET
		address_line_2 = $1
		WHERE address_id = $2`,
		line,
		in.AddressId,
	)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "user.address.104",
			Interface: err,
		})
	}

	rb.BuildResponse(ctx, out)
}

var AddressValidationError = uf.ApiError{
	Code: "address_validation_error",
	Msg:  "USPS address validation failed",
}

func HAddressValidateUsps(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	uf.Debug("HAddressValidateUsps coming in")
	in, out := AddressIdStruct{}, uf.IsValidStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))

	address, err := LAddressGetById(&gibs, in.AddressId)
	if err != nil {
		uf.Glog(&gibs, uf.GlogStruct{
			Level:     uf.LevelError,
			Code:      "user.address.105",
			Interface: err,
		})
		rb.AddError(GetAddressError)
	}

	uf.Trace("found address in db")
	uf.Trace(in)
	uf.Trace(address)

	fresh := UspsAddressValidation(&gibs, address)
	uf.Glog(&gibs, uf.GlogStruct{
		Level:     uf.LevelTrace,
		Code:      "user.address.900",
		Interface: err,
	})
	if fresh.Address.City == "" {
		out.IsValid = false

	} else {
		//NOTE: let's leave it to the user to type their apartment correctly
		address.AddressLine1 = fresh.Address.Address2
		address.City = fresh.Address.City
		address.State = fresh.Address.State
		address.Postal = fresh.Address.Zip5
		address.Zip4 = fresh.Address.Zip4

		address.IsBuilder = false
		address.IsValidated = true
		address.LAddressUpsert(&gibs)

		out.IsValid = true
	}

	rb.BuildResponse(ctx, out)
}
