package uf

type ApiError struct {
	Code string
	Msg  string
}

func CustomError(err error) ApiError {
	return ApiError{
		Code: err.Error(),
		Msg:  err.Error(),
	}
}

var PolicyDeniedError = ApiError{
	Code: "policy_denied",
	Msg:  "policy denied",
}
var PolicyNotInspectedError = ApiError{
	Code: "policy_not_inspected",
	Msg:  "policy not inspected",
}
var RequirementsMissingError = ApiError{
	Code: "requirements_missing",
	Msg:  "requirements missing",
}

var EncodeUfRequestBodyError = ApiError{
	Code: "encode_uf_request_body_error",
	Msg:  "Failed to encode UfRequest.Body.",
}
var UfRequestError = ApiError{
	Code: "uf_request_error",
	Msg:  "UfRequest failed.",
}
var DecodeUfResponseBodyError = ApiError{
	Code: "decode_uf_response_body_error",
	Msg:  "Failed to decode UfResponse.Body.",
}
var DecodeUfResponseErrorsError = ApiError{
	Code: "decode_uf_response_errors_error",
	Msg:  "Failed to decode UfResponse.Body.Errors.",
}
