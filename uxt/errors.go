package uxt

import "gitlab.com/textfridayy/uno/uf"

var SurferNotFoundError = uf.ApiError{
	Code: "surfer_not_found",
	Msg:  "Surfer not found",
}

var FailedToCreateJwtError = uf.ApiError{
	Code: "failed_to_create_jwt",
	Msg:  "Failed to login",
}
