package uf_order

import (
	"net/http"

	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/zinc"
)

func GetOrderInfo(
	conf *uf.Config,
	zincRequestId string) (zinc.GetOrderOut, *http.Response, error) {

	path := zinc.PathOrders + "/" + zincRequestId
	info := zinc.GetOrderOut{}
	httpResp, err := zinc.DoRequest(
		conf,
		path,
		uf.HttpMethodGet,
		nil,
		nil,
		&info,
	)

	if info.Error.Code != "" {
		uf.Fatal(info.Error.Code + "\n" + info.Error.Data.Message)
	}

	return info, httpResp, err
}
