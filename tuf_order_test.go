package main

import (
	"io"
	"testing"

	"gitlab.com/textfridayy/uno/scaffold"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_order"
	"gitlab.com/textfridayy/uno/zinc"
)

func TestExecuteCreateOrder(t *testing.T) {
	pile := uf.Pile{Conf: uf.Config{}}
	pile.Conf.InterpretDefaults("order")

	shippingInfo := zinc.ShippingObject{
		OrderBy:  "price",
		MaxPrice: 0,
	}

	products := []zinc.ProductObject{scaffold.ProductObject01()}

	order, orderResp, httpResp, err := uf_order.ExecuteCreateOrder(
		&pile.Conf,
		scaffold.SurferPg01(),
		scaffold.AddressObject01(),
		scaffold.AddressObject01(),
		shippingInfo,
		products,
		0,
	)
	if err != nil {
		t.Fatalf("Encountered error: \n%v", err)
	}

	uf.Debug(order.ShippingAddress.FirstName)
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		t.Fatalf("Encountered error: \n%v", err)
	}
	uf.Debug(string(body))
	uf.Debug(httpResp.StatusCode)
	uf.Debug(orderResp)
}

//func TestGetOrderDetails(t *testing.T) {
//pile := uf.Pile{Conf: uf.Config{}}
//pile.Conf.InterpretDefaults("order")

//info, httpResp, err := uf_order.GetOrderInfo(
//&pile.Conf,
//"0c4c7762b400a5bef24003f864fb6ead",
//)
////"df85a1a67e02100b94dce4e14fa0dc3d",
//if err != nil {
//t.Fatalf("Encountered error: \n%v", err)
//}

//uf.Debug(info)
//uf.Debug(info.Error.Code)
//body, err := io.ReadAll(httpResp.Body)
//if err != nil {
//t.Fatalf("Encountered error: \n%v", err)
//}
//uf.Debug(string(body))
//uf.Debug(httpResp.StatusCode)

//if info.Error.Code != "" {
//uf.Fatal(info.Error.Code + "\n" + info.Error.Data.Message)
//}
//}

//func TestRetryFailedOrder(t *testing.T) {
//pile := uf.Pile{Conf: uf.Config{}}
//pile.Conf.InterpretDefaults("order")

//info, httpResp, err := uf_order.RetryFailedOrder(
//&pile.Conf,
//"0c4c7762b400a5bef24003f864fb6ead",
//)
//if err != nil {
//t.Fatalf("Encountered error: \n%v", err)
//}

//uf.Debug(info)
//body, err := io.ReadAll(httpResp.Body)
//if err != nil {
//t.Fatalf("Encountered error: \n%v", err)
//}
//uf.Debug(string(body))
//uf.Debug(httpResp.StatusCode)
//}
