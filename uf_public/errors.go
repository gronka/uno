package uf_public

import "gitlab.com/textfridayy/uno/uf"

var CredentialsError = uf.ApiError{
	Code: "credentials_error",
	Msg:  "Failed to sign in.",
}
