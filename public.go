package validator

import (
	"github.com/joaosoft/errors"
)

func AddBefore(name string, handler BeforeTagHandler) *Validator {
	return validator.AddBefore(name, handler)
}

func AddMiddle(name string, handler MiddleTagHandler) *Validator {
	return validator.AddMiddle(name, handler)
}

func AddAfter(name string, handler AfterTagHandler) *Validator {
	return validator.AddAfter(name, handler)
}

func SetValidateAll(validate bool) {
	validator.SetValidateAll(validate)
}

func SetTag(tag string) {
	validator.SetTag(tag)
}

// Validate ...
func Validate(obj interface{}) *errors.ListErr {
	return validator.Validate(obj)
}
