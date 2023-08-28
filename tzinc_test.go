package main

import (
	"os"
	"testing"

	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/zinc"
)

func TestProductSearch(t *testing.T) {
	pile := uf.Pile{Conf: uf.Config{}}
	os.Setenv("Environment", "local")
	pile.Conf.InterpretDefaults("border")

	srb, err := zinc.SearchForProduct(&pile.Conf, "soap", "amazon")
	if err != nil {
		t.Fatalf("Encountered error: \n%v", err)
	}

	if len(srb.Results) < 3 {
		t.Fatalf("Not enough results found in response: \n%v", srb)

	}

	uf.Debug(srb.PPrintFirst())
}

//func TestGetProductDetails(t *testing.T) {
//pile := uf.Pile{Conf: uf.Config{}}
//pile.Conf.InterpretDefaults("border")

//srb, err := zinc.GetProductDetails(&pile.Conf, "B079FV6PRH", "amazon")
//if err != nil {
//t.Fatalf("Encountered error: \n%v", err)
//}

//uf.Debug(srb)
//}
