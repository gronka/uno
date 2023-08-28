package scaffold

import (
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
	"gitlab.com/textfridayy/uno/zinc"
)

func SurferPg01() uf_user.SurferPg {
	return uf_user.SurferPg{
		SurferId:                   uf.ZerosUuid,
		Name:                       "Taylor Gronka",
		Email:                      "mr.gronka@gmail.com",
		Phone:                      "+12345678912",
		StripeCustomerId:           "testing_id",
		StripeDefaultPaymentMethod: "method2",
	}
}

func AddressObject01() zinc.AddressObject {
	return zinc.AddressObject{
		FirstName:    "Jane",
		LastName:     "Doe",
		AddressLine1: "1924 Dauphin Island Pkwy B",
		AddressLine2: "",
		ZipCode:      "36605",
		City:         "Mobile",
		State:        "Alabama",
		Country:      "USA",
		PhoneNumber:  "+12345678912",
		Instructions: "You can leave it on the steps",
	}
}

func ProductObjectWithId(id string) zinc.ProductObject {
	return zinc.ProductObject{
		ProductId: "B079FV6PRH",
		Quantity:  1,
		SellerSelectionCriteria: zinc.SellerSelectionCriteriaObject{
			ConditionIn:     []string{"new"},
			HandlingDaysMax: 10,
		},
	}
}

func ProductObject01() zinc.ProductObject {
	return zinc.ProductObject{
		ProductId: "B079FV6PRH",
		Quantity:  1,
		SellerSelectionCriteria: zinc.SellerSelectionCriteriaObject{
			ConditionIn:     []string{"new"},
			HandlingDaysMax: 10,
		},
	}
}
