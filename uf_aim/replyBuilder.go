package uf_aim

type ReplyBuilder struct {
	Body   string
	Errors []ReplyError
}

type ReplyError struct {
	Code string
	Msg  string
}
