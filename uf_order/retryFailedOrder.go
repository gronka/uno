package uf_order

import (
	"net/http"

	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/zinc"
)

func RetryFailedOrder(
	conf *uf.Config,
	zincRequestId string) (zinc.RequestIdBody, *http.Response, error) {

	path := zinc.PathOrders + "/" + zincRequestId + "/retry"
	out := zinc.RequestIdBody{}
	httpResp, err := zinc.DoRequest(
		conf,
		path,
		uf.HttpMethodPost,
		nil,
		nil,
		&out,
	)

	return out, httpResp, err
}
