package uf_order

import (
	"net/http"

	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/zinc"
)

// An order can be aborted while Zinc is still processing it. After zinc
// processed it, it can only be Cancelled if placed on Amazon. If not placed
// on Amazon, the customer must contact the seller directly.
func AbortOrTryCancelOrder(
	conf *uf.Config,
	zincRequestId string) (zinc.GetOrderOut, *http.Response, error) {

	info, httpResp, err := GetOrderInfo(
		conf,
		zincRequestId,
	)
	if err != nil {
		uf.Fatal(err)
	}

	//TODO: what is state of failed or completed order?
	//TODO: what is state of cancelled order?
	//TODO: ask zinc if the source code is available

	switch info.Error.Code {
	case "aborted_request":
		//TODO reply already cancelled
	case "request_processing":
		//info, httpResp, err := AbortOrder(conf, zincRequestId)
		abortResponse, _, err := AbortOrder(conf, zincRequestId)
		if err != nil {
			//TODO error contacting zinc
		}
		if abortResponse.Code == "aborted_request" {
			//TODO reply already cancelled
		} else if abortResponse.Code == "request_processing" {
			//TODO reply abort failed (why would this happen??)
		}
	case "attempting_to_cancel":
		//TODO amazon is trying to cancel this order
	default:
		if info.MerchantOrderIds[0].Merchant != "amazon" {
			//TODO to cancel this order, you must contact the seller directly. Here
			// is their information
		} else {
			//TODO for order in range
			for _, merchantObject := range info.MerchantOrderIds {
				cancelResponse, _, err := TryCancelOrder(
					conf,
					zincRequestId,
					merchantObject.MerchantOrderId)
				uf.Debug(cancelResponse.ZincCancelRequestId)
				uf.Debug(err)

				//TODO handle if err != nil
			}
			//TODO: save cancelResponse.RequestId
		}
	}

	return info, httpResp, err
}

func AbortOrder(
	conf *uf.Config,
	zincRequestId string) (zinc.OrderAbortResponse, *http.Response, error) {

	path := zinc.PathOrders + "/" + zincRequestId + "/abort"
	response := zinc.OrderAbortResponse{}
	httpResp, err := zinc.DoRequest(
		conf,
		path,
		uf.HttpMethodPost,
		nil,
		nil,
		&response,
	)

	if response.Code != "" {
		uf.Fatal(response.Code + "\n" + response.Message)
	}

	return response, httpResp, err
}

func TryCancelOrder(
	conf *uf.Config,
	zincRequestId string,
	merchantOrderId string,
) (zinc.TryOrderCancelResponse, *http.Response, error) {

	path := zinc.PathOrders + "/" + zincRequestId + "/cancel"
	body := zinc.TryOrderCancelRequest{
		MerchantOrderId: merchantOrderId,
		Webhooks: zinc.WebhooksObject{
			RequestSucceeded: conf.OrderAddress + "/zinc.cancelOrderSuccess",
			RequestFailed:    conf.OrderAddress + "/zinc.cancelOrderFail",
		},
		//TODO does this support ClientNotesMeta? it's used in Fridayy2
	}

	response := zinc.TryOrderCancelResponse{}
	httpResp, err := zinc.DoRequest(
		conf,
		path,
		uf.HttpMethodPost,
		body,
		nil,
		&response,
	)

	return response, httpResp, err
}
