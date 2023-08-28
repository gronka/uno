package zinc

import "gitlab.com/textfridayy/uno/uf"

type OrderAbortResponse struct {
	Type    string `json:"_type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	// Data to be the same as Message
	Data      map[string]interface{} `json:"data"`
	RequestId string                 `json:"request_id"`
	// original request for order creation
	Request map[string]interface{} `json:"request"`
}

type TryOrderCancelRequest struct {
	MerchantOrderId string         `json:"merchant_order_id"`
	Webhooks        WebhooksObject `json:"webhooks"`
}

type TryOrderCancelResponse struct {
	ZincCancelRequestId string `json:"request_id"`
}

type CancellationInfoResponse struct {
	Type            string                 `json:"_type"`
	MerchantOrderId string                 `json:"merchant_order_id"`
	Request         map[string]interface{} `json:"request"`
}

func OrderCancel(
	conf *uf.Config,
	orderRequestId string) (*GetOrderOut, error) {
	uf.Trace("zinc get order details for requestId: " + orderRequestId)

	out := GetOrderOut{}
	_, err := DoRequest(
		conf,
		string(PathOrders)+"/"+orderRequestId+"/abort",
		uf.HttpMethodPost,
		nil,
		nil,
		&out,
	)
	uf.Trace("returning results")

	return &out, err
}
