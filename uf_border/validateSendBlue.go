package uf_border

import "gitlab.com/textfridayy/uno/uf"

const SmsContentMaxLength = 120
const SmsContentMinLength = 0

type Validation struct {
	uf.Validation
}

func (in *SendBlueIn) ValidateAllFields() (val Validation) {
	val.SmsPhone(in.FromNumber)
	val.SmsContent(in.Content)
	return
}

var CountryNotServicedError = uf.ApiError{
	Code: "country_not_serviced",
	Msg:  "Sorry - Fridayy does not yet operate in your country.",
}

var UnrecognizedNumberFormatError = uf.ApiError{
	Code: "unrecognized_number_format",
	Msg:  "Sorry - Fridayy cannot recognize your number.",
}

func (val *Validation) SmsPhone(phone string) *Validation {
	if phone[0:2] != "+1" {
		val.AddError(CountryNotServicedError)
	}

	if len(phone) != 12 {
		val.AddError(UnrecognizedNumberFormatError)
	}
	return val
}

var SmsContentTooShortError = uf.ApiError{
	Code: "sms_content_too_short",
	Msg:  "Message is too short to process",
}

func (val *Validation) SmsContent(content string) *Validation {
	if len(content) <= 1 {
		val.AddError(SmsContentTooShortError)
	}
	return val
}
