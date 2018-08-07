package validator

import (
	"github.com/joaosoft/errors"
)

func AddPre(name string, handler PreTagHandler) *Validator {
	return validator.AddPre(name, handler)
}

func AddMiddle(name string, handler MiddleTagHandler) *Validator {
	return validator.AddMiddle(name, handler)
}

func AddPos(name string, handler PosTagHandler) *Validator {
	return validator.AddPos(name, handler)
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
