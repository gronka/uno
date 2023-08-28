package ut

import (
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/zinc"
)

type OrderDetails struct {
	Tracking []TrackingDetails
}

type TrackingDetails struct {
	// Carrier such as "fedex"
	Carrier string
	// DeliveryStatus is "delivered" after delivery
	DeliveryStatus      string
	TrackingNumber      string
	TrackingUrl         string
	ZincProductIds      []string
	TimeTrackingUpdated int64
}

func GetOrderDetailsByDriver(
	gibs *uf.Gibs,
	order *OrderPg) (od OrderDetails, err error) {

	od.Tracking = make([]TrackingDetails, 0)

	switch order.Driver {
	case "zinc":
		zincDetails, err := zinc.GetOrderDetails(gibs.Conf, order.ZincOrderRequestId)
		if err != nil {
			return od, err
		}

		for _, zincTracking := range zincDetails.Tracking {
			//TODO: convert obtainedAt to a timestamp
			// obtainedAt format: 2018-11-29T20:53:05.924Z
			obtainedAt := zincTracking.ObtainedAt
			uf.Trace(obtainedAt)

			od.Tracking = append(od.Tracking, TrackingDetails{
				Carrier:        zincTracking.Carrier,
				DeliveryStatus: zincTracking.DeliveryStatus,
				TrackingNumber: zincTracking.TrackingNumber,
				TrackingUrl:    zincTracking.TrackingUrl,
				ZincProductIds: zincTracking.ProductIds,
			})

		}

	default:
		uf.Glog(gibs, uf.GlogStruct{
			Level: uf.LevelError,
			Code:  "order.100",
			Msg:   "Invalid order driver: " + order.Driver,
		})
	}
	return od, err
}

type CancelDetails struct {
	Status string
}

func OrderCancelByDriver(
	gibs *uf.Gibs,
	order *OrderPg) (cd CancelDetails, err error) {

	switch order.Driver {
	case "zinc":
		orderObject, err := zinc.OrderCancel(gibs.Conf, order.ZincOrderRequestId)
		if err != nil {
			return cd, err
		}

		switch orderObject.Code {
		case "aborted_request":
			cd.Status = "success"
		case "request_processing":
			cd.Status = "attempting"
		default:
			uf.Glog(gibs, uf.GlogStruct{
				Level: uf.LevelError,
				Code:  "order.102",
				Msg:   "Unknown zinc order status: " + orderObject.Code,
			})
		}

	default:
		uf.Glog(gibs, uf.GlogStruct{
			Level: uf.LevelError,
			Code:  "order.103",
			Msg:   "Invalid order driver: " + order.Driver,
		})
	}
	return cd, err
}
