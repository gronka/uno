package uxt

import (
	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
)

func AddressCreateBuilder(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) (uf_user.AddressPg, uf.UfResponse) {
	pkg := uf_user.SurferIdStruct{SurferId: surferId}

	address := uf_user.AddressPg{}
	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressCreateBuilderV1,
		pkg,
		&address,
	)

	return address, ures
}

func AddressDelete(
	gibs *uf.Gibs,
	addressId gocql.UUID,
) uf.UfResponse {
	pkg := uf_user.AddressIdStruct{AddressId: addressId}

	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressDeleteV1,
		pkg,
		uf.EmptyStruct{},
	)

	return ures
}

func AddressDeleteBuilder(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) uf.UfResponse {
	pkg := uf_user.SurferIdStruct{SurferId: surferId}

	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressDeleteBuilderV1,
		pkg,
		uf.EmptyStruct{},
	)

	return ures
}

func AddressGetById(
	gibs *uf.Gibs,
	addressId gocql.UUID,
) (uf_user.AddressPg, uf.UfResponse) {
	pkg := uf_user.AddressIdStruct{AddressId: addressId}

	address := uf_user.AddressPg{}
	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressGetByIdV1,
		pkg,
		&address,
	)

	return address, ures
}

func AddressGetBuilderBySurferId(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) (uf_user.AddressPg, uf.UfResponse) {
	pkg := uf_user.SurferIdStruct{SurferId: surferId}

	address := uf_user.AddressPg{}
	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressGetBuilderBySurferIdV1,
		pkg,
		&address,
	)

	return address, ures
}

func AddressGetNonBuilderBySurferId(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) (uf_user.AddressPg, uf.UfResponse) {
	pkg := uf_user.SurferIdStruct{SurferId: surferId}

	address := uf_user.AddressPg{}
	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressGetNonBuilderBySurferIdV1,
		pkg,
		&address,
	)

	return address, ures
}

func AddressUpdateName(
	gibs *uf.Gibs,
	addressId gocql.UUID,
	name string,
) uf.UfResponse {
	pkg := uf_user.AddressUpdateNameStruct{
		AddressId: addressId,
		Name:      name,
	}

	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressUpdateNameV1,
		pkg,
		uf.EmptyStruct{},
	)

	return ures
}

func AddressUpdatePostalPlus(
	gibs *uf.Gibs,
	addressId gocql.UUID,
	postal string,
) uf.UfResponse {
	pkg := uf_user.AddressUpdatePostalStruct{
		AddressId: addressId,
		Postal:    postal,
	}

	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressUpdatePostalPlusV1,
		pkg,
		uf.EmptyStruct{},
	)

	return ures
}

func AddressUpdateLine1(
	gibs *uf.Gibs,
	addressId gocql.UUID,
	line1 string,
) uf.UfResponse {
	pkg := uf_user.AddressUpdateAddressLine1Struct{
		AddressId:    addressId,
		AddressLine1: line1,
	}

	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressUpdateLine1V1,
		pkg,
		uf.EmptyStruct{},
	)

	return ures
}

func AddressUpdateLine2(
	gibs *uf.Gibs,
	addressId gocql.UUID,
	line2 string,
) uf.UfResponse {
	pkg := uf_user.AddressUpdateAddressLine2Struct{
		AddressId:    addressId,
		AddressLine2: line2,
	}

	ures := uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressUpdateLine2V1,
		pkg,
		uf.EmptyStruct{},
	)

	return ures
}

func AddressValidateUsps(
	gibs *uf.Gibs,
	addressId gocql.UUID,
) (valid uf.IsValidStruct, ures uf.UfResponse) {
	pkg := uf_user.AddressIdStruct{
		AddressId: addressId,
	}

	ures = uf.MakeRequest(
		gibs,
		gibs.Conf.UserAddress,
		UserAddressValidateUspsV1,
		pkg,
		&valid,
	)

	return valid, ures
}
