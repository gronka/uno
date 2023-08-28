package uf_border

import "gitlab.com/textfridayy/uno/uf"

type JsonIn struct {
	Phone   string
	Content string
}

const ContentMaxLength = 120
const ContentMinLength = 0

func (in *JsonIn) ValidateAllFields() (val Validation) {
	val.JsonPhone(in.Phone)
	val.JsonContent(in.Content)
	return
}

func (val *Validation) JsonPhone(phone string) *Validation {
	if phone[0:2] != "+1" {
		val.AddError(CountryNotServicedError)
	}

	if len(phone) != 12 {
		val.AddError(UnrecognizedNumberFormatError)
	}
	return val
}

var JsonContentTooShortError = uf.ApiError{
	Code: "json_content_too_short",
	Msg:  "Message is too short to process",
}

func (val *Validation) JsonContent(content string) *Validation {
	if len(content) <= 1 {
		val.AddError(JsonContentTooShortError)
	}
	return val
}
