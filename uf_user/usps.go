package uf_user

import (
	"fmt"

	usps "github.com/AbhiAgarwal/usps-go"
	"gitlab.com/textfridayy/uno/uf"
)

//NOTE we could also use melissa.com or algolia.com for verification
func UspsZipcodeToCityState(gibs *uf.Gibs, zip string) usps.CityStateLookupResponse {
	var uspsConn usps.USPS
	uspsConn.Username = gibs.Conf.UspsUsername
	uspsConn.Password = gibs.Conf.UspsPassword

	var address usps.ZipCode
	address.Zip5 = zip

	output := uspsConn.CityStateLookup(address)

	trace := fmt.Sprintf("zipcode %v found %v, %v",
		zip,
		output.ZipC.City,
		output.ZipC.State,
	)
	uf.Trace(trace)
	return output
}

func UspsAddressValidation(gibs *uf.Gibs, address AddressPg) usps.AddressValidateResponse {
	var uspsConn usps.USPS
	uspsConn.Username = gibs.Conf.UspsUsername
	uspsConn.Password = gibs.Conf.UspsPassword

	var verify usps.Address
	//NOTE: USPS swaps lines 1 and 2 for some reason
	verify.Address2 = address.AddressLine1
	//NOTE: let's leave it to the user to type their apartment correctly
	//verify.Address1 = address.AddressLine2
	verify.City = address.City
	verify.State = address.State

	output := uspsConn.AddressVerification(verify)

	trace := fmt.Sprintf("validation %v, %v, %v",
		output.Address.Address2,
		output.Address.City,
		output.Address.State,
	)
	uf.Trace(trace)
	return output
}
