package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/joaosoft/errors"
	"github.com/satori/go.uuid"
)

func (v *Validator) validate_value(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.ListErr {
	rtnErrs := make(errors.ListErr, 0)

	if fmt.Sprintf("%+v", value) == "" || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return rtnErrs
	}

	if fmt.Sprintf("%+v", value) != fmt.Sprintf("%+v", expected) {
		rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("the value [%+v] is different of the expected [%+v] on field [%s]", value, expected, name)))
	}

	return rtnErrs
}

func (v *Validator) validate_options(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.ListErr {
	rtnErrs := make(errors.ListErr, 0)

	if fmt.Sprintf("%+v", value) == "" || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return rtnErrs
	}

	options := strings.Split(expected.(string), ";")
	var invalidValue interface{}

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		optionsVal := make(map[string]bool)
		for _, option := range options {
			optionsVal[option] = true
		}

		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)

			_, ok := optionsVal[fmt.Sprintf("%+v", nextValue.Interface())]
			if !ok {
				invalidValue = nextValue.Interface()
				rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("the value [%+v] is different of the expected options [%+v] on field [%s]", invalidValue, expected, name)))
				if !v.validateAll {
					break
				}
			}
		}

	case reflect.Map:
		optionsMap := make(map[string]interface{})
		for _, option := range options {
			values := strings.Split(option, ":")
			if len(values) != 2 {
				continue
			}
			optionsMap[values[0]] = values[1]
		}

		for _, key := range value.MapKeys() {
			nextValue := value.MapIndex(key)

			val, ok := optionsMap[fmt.Sprintf("%+v", key.Interface())]
			if !ok || fmt.Sprintf("%+v", nextValue.Interface()) != fmt.Sprintf("%+v", val) {
				invalidValue = fmt.Sprintf("%+v:%+v", key.Interface(), nextValue.Interface())
				e := errors.New("0", fmt.Sprintf("the value [%+v] is different of the expected options [%+v] on field [%s]", nextValue.Interface(), expected, name))
				rtnErrs = append(rtnErrs, e)
				if !v.validateAll {
					break
				}
			}
		}

	default:
		optionsVal := make(map[string]bool)
		for _, option := range options {
			optionsVal[option] = true
		}

		_, ok := optionsVal[fmt.Sprintf("%+v", value)]
		if !ok {
			invalidValue = value
			rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("the value [%+v] is different of the expected options [%+v] on field [%s]", invalidValue, expected, name)))
		}
	}

	return rtnErrs
}

func (v *Validator) validate_size(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.ListErr {
	rtnErrs := make(errors.ListErr, 0)
	size, e := strconv.Atoi(expected.(string))
	if e != nil {
		rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("the size [%s] is invalid on field [%s]", expected, value)))
		return rtnErrs
	}

	var valueSize int64

	switch value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		valueSize = int64(value.Len())
	case reflect.String:
		valueSize = int64(len(value.String()))
	default:
		valueSize = value.Int()
	}

	if valueSize != int64(size) {
		rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("the length [%+v] is lower then the expected [%+v] on field [%s]", valueSize, expected, name)))
	}

	return rtnErrs
}

func (v *Validator) validate_min(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.ListErr {
	rtnErrs := make(errors.ListErr, 0)
	min, e := strconv.Atoi(expected.(string))
	if e != nil {
		rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("the size [%s] is invalid on field [%s]", expected, value)))
		return rtnErrs
	}

	var valueSize int64

	switch value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		valueSize = int64(value.Len())
	case reflect.String:
		valueSize = int64(len(value.String()))
	default:
		valueSize = value.Int()
	}

	if valueSize < int64(min) {
		rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("the length [%+v] is lower then the expected [%+v] on field [%s]", valueSize, expected, name)))
	}

	return rtnErrs
}

func (v *Validator) validate_max(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.ListErr {
	rtnErrs := make(errors.ListErr, 0)
	max, e := strconv.Atoi(expected.(string))
	if e != nil {
		rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("the size [%s] is invalid on field [%s]", expected, value)))
		return rtnErrs
	}

	var valueSize int64

	switch value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		valueSize = int64(value.Len())
	case reflect.String:
		valueSize = int64(len(value.String()))
	default:
		valueSize = value.Int()
	}

	if valueSize > int64(max) {
		rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("the length [%+v] is bigger then the expected [%+v] on field [%s]", valueSize, expected, name)))
	}

	return rtnErrs
}

func (v *Validator) validate_nonzero(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.ListErr {
	rtnErrs := make(errors.ListErr, 0)
	var valueSize int64

	switch value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:

		switch value.Type() {
		case reflect.TypeOf(uuid.UUID{}):
			if value.Interface().(uuid.UUID) != uuid.Nil {
				valueSize = 1
			}
		default:
			valueSize = int64(value.Len())
		}

	case reflect.String:
		valueSize = int64(len(value.String()))

	default:
		valueSize = value.Int()
	}

	if valueSize == 0 {
		rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("the value shouldn't be zero on field [%s]", name)))
	}

	return rtnErrs
}

func (v *Validator) validate_regex(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.ListErr {

	rtnErrs := make(errors.ListErr, 0)

	if fmt.Sprintf("%+v", value) == "" || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return rtnErrs
	}

	r, e := regexp.Compile(expected.(string))
	if e != nil {
		rtnErrs = append(rtnErrs, errors.New("0", e))
		return rtnErrs
	}

	if len(fmt.Sprintf("%+v", value)) > 0 {
		if !r.MatchString(fmt.Sprintf("%+v", value)) {
			rtnErrs = append(rtnErrs, errors.New("0", fmt.Sprintf("invalid data [%s] on field [%+v] ", value, name)))
		}
	}

	return rtnErrs
}

func (v *Validator) validate_special(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.ListErr {

	if fmt.Sprintf("%+v", value) == "" || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return make(errors.ListErr, 0)
	}

	switch expected {
	case RegexTagForDateDefault:
		expected = RegexForDateDefault
	case RegexTagForDateDDMMYYYY:
		expected = RegexForDateDDMMYYYY
	case RegexTagForDateYYYYMMDD:
		expected = RegexForDateYYYYMMDD
	case RegexTagForTimeDefault:
		expected = RegexForTimeDefault
	case RegexTagForTimeHHMMSS:
		expected = RegexForTimeHHMMSS
	default:
		return []*errors.Err{errors.New("0", fmt.Sprintf("invalid special [%s] on field [%+v] ", expected, name))}
	}

	return v.validate_regex(name, value, expected, errs)
}

func (v *Validator) validate_error(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.ListErr {
	rtnErrs := make(errors.ListErr, 0)
	added := make(map[string]bool)
	for i, _ := range *errs {
		(*errs)[i].SetCode(expected.(string))

		if v.errorCodeHandler != nil {
			if matched, err := regexp.MatchString("{[a-z0-9]+}", expected.(string)); err != nil {
				rtnErrs = append(rtnErrs, errors.New("0", err))
			} else {
				if matched {
					replacer := strings.NewReplacer("{", "", "}", "")
					errorCode := replacer.Replace(expected.(string))

					if _, ok := added[errorCode]; !ok {
						newErr := v.errorCodeHandler(errorCode)
						(*errs)[i].SetErr(newErr.(*errors.Err))

						added[errorCode] = true
					} else {
						*errs = append((*errs)[:i], (*errs)[i+1:]...)
					}
				}
			}
		}
	}

	return rtnErrs
}
