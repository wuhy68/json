package validator

import (
	"reflect"

	"regexp"

	"strings"

	"github.com/joaosoft/errors"
)

func (v *Validator) NewDefaultPosHandlers() map[string]PosTagHandler {
	return map[string]PosTagHandler{"error": v.pos_error}
}

func (v *Validator) pos_error(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.IErr {

	for i, _ := range *errs {
		(*errs)[i].SetCode(expected.(string))

		if v.errorCodeHandler != nil {
			if matched, err := regexp.MatchString("{[a-z0-9]+}", expected.(string)); err != nil {
				return errors.New("0", err)
			} else {
				if matched {
					replacer := strings.NewReplacer("{", "", "}", "")

					errorCode := replacer.Replace(expected.(string))
					newErr := v.errorCodeHandler(errorCode)
					(*errs)[i].SetError(newErr.(*errors.Err))
				}
			}
		}
	}

	return nil
}
