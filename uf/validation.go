package uf

type Validation struct {
	Errors []ApiError
}

func (val *Validation) IsValid() bool {
	if val.Errors != nil && len(val.Errors) > 0 {
		return false
	}
	return true
}
func (val *Validation) IsNotValid() bool {
	return !val.IsValid()
}

func (val *Validation) AddError(valError ApiError) {
	if val.Errors == nil {
		val.Errors = []ApiError{valError}
	} else {
		val.Errors = append(val.Errors, valError)
	}
}

// Perform, in an abstract sense, tells a Validation object to perform a
// validation, which is passed as another Validation object. The result is the
// same as combining two Validation objects.
//func (validation *Validation) Perform(val *Validation) {
//if validation.Errors == nil {
//validation.Errors = val.Errors
//} else {
//validation.Errors = append(validation.Errors, val.Errors...)
//}
//}
