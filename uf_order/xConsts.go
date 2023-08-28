package uf_order

import "gitlab.com/textfridayy/uno/zinc"

func getTylerBillingAddressObject() zinc.AddressObject {
	return zinc.AddressObject{
		FirstName:    "Tyler",
		LastName:     "Mosher",
		AddressLine1: "42 Labor In Vain",
		AddressLine2: "",
		ZipCode:      "01938",
		City:         "Ipswich",
		State:        "MA",
		Country:      "US",
		//TODO: what phone number to use - this one? ADMIN_NUMBER = '+16502792417'
		PhoneNumber:  "2678888227",
		Instructions: "",
	}
}
