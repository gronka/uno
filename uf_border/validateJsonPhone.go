package uf_border

func (in *JsonPhoneIn) ValidateAllFields() (val Validation) {
	val.JsonPhone(in.Phone)
	val.JsonContent(in.Content)
	return
}

func (val *Validation) JsonPhonePhone(phone string) *Validation {
	if phone[0:1] != "1" {
		val.AddError(CountryNotServicedError)
	}

	//TODO use a phone validation library
	if len(phone) != 11 {
		val.AddError(UnrecognizedNumberFormatError)
	}
	return val
}

func (val *Validation) JsonPhoneContent(content string) *Validation {
	if len(content) <= 1 {
		val.AddError(JsonContentTooShortError)
	}
	return val
}
