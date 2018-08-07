package validator

import (
	"reflect"

	"github.com/joaosoft/errors"
)

func NewDefaultPosHandlers() map[string]PosTagHandler {
	return map[string]PosTagHandler{"error": pos_error}
}

func pos_error(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) *errors.Err {

	for i, _ := range *errs {
		(*errs)[i].Code = expected.(string)
	}

	return nil
}
