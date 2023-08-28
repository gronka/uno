package uf_user

import "gitlab.com/textfridayy/uno/uf"

var GetUserError = uf.ApiError{
	Code: "get_user_error",
	Msg:  "Error while getting user.",
}

var CreateUserError = uf.ApiError{
	Code: "create_user_error",
	Msg:  "Error while creating user.",
}
