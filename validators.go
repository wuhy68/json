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

func (v *Validator) validate_value(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
	if fmt.Sprintf("%+v", value) != fmt.Sprintf("%+v", expected) {
		return errors.New("0", fmt.Sprintf("the value [%+v] is different of the expected [%+v] on field [%s]", value, expected, name))
	}

	return nil
}

func (v *Validator) validate_options(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
	options := strings.Split(expected.(string), ";")
	var valid = true
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
				valid = false
				break
			}
		}

	case reflect.Map:
		valid = true
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
				valid = false
				break
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
			valid = false
			break
		}
	}

	if !valid {
		return errors.New("0", fmt.Sprintf("the value [%+v] is different of the expected options [%+v] on field [%s]", invalidValue, expected, name))
	}

	return nil
}

func (v *Validator) validate_size(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
	size, e := strconv.Atoi(expected.(string))
	if e != nil {
		return errors.New("0", fmt.Sprintf("the size [%s] is invalid on field [%s]", expected, value))
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
		return errors.New("0", fmt.Sprintf("the length [%+v] is lower then the expected [%+v] on field [%s]", valueSize, expected, name))
	}

	return nil
}

func (v *Validator) validate_min(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
	min, e := strconv.Atoi(expected.(string))
	if e != nil {
		return errors.New("0", fmt.Sprintf("the size [%s] is invalid on field [%s]", expected, value))
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
		return errors.New("0", fmt.Sprintf("the length [%+v] is lower then the expected [%+v] on field [%s]", valueSize, expected, name))
	}

	return nil
}

func (v *Validator) validate_max(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
	max, e := strconv.Atoi(expected.(string))
	if e != nil {
		return errors.New("0", fmt.Sprintf("the size [%s] is invalid on field [%s]", expected, value))
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
		return errors.New("0", fmt.Sprintf("the length [%+v] is bigger then the expected [%+v] on field [%s]", valueSize, expected, name))
	}

	return nil
}

func (v *Validator) validate_nonzero(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
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
		return errors.New("0", fmt.Sprintf("the value shouldn't be zero on field [%s]", name))
	}

	return nil
}

func (v *Validator) validate_regex(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {

	r, e := regexp.Compile(expected.(string))
	if e != nil {
		return errors.New("0", e)
	}

	if !r.MatchString(fmt.Sprintf("%+v", value)) {
		return errors.New("0", fmt.Sprintf("invalid data [%s] on field [%+v] ", value, name))
	}

	return nil
}
