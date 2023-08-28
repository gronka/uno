package zinc

type ErrorObject struct {
	Type string          `json:"_type"`
	Code string          `json:"code"`
	Data ErrorDataObject `json:"data"`
}

type ErrorDataObject struct {
	Details       string `json:"details"`
	Message       string `json:"message"`
	CurrentUrl    string `json:"current_url"`
	OrderResponse string `json:"order_response"`
}
