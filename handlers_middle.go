package validator

import (
	"fmt"

	"reflect"

	"strconv"

	"regexp"

	"github.com/joaosoft/errors"
	"github.com/satori/go.uuid"
)

func NewDefaultMiddleHandlers() map[string]MiddleTagHandler {
	return map[string]MiddleTagHandler{"value": middle_value, "size": middle_size, "min": middle_min, "max": middle_max, "nonzero": middle_nonzero, "regex": middle_regex}
}

func middle_value(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
	if fmt.Sprintf("%+v", value) != fmt.Sprintf("%+v", expected) {
		return errors.New("0", fmt.Sprintf("the value [%+v] is diferent of the expected [%+v] on field [%s]", value, expected, name))
	}

	return nil
}

func middle_size(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
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

func middle_min(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
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

func middle_max(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
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

func middle_nonzero(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
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
		return errors.New("0", fmt.Sprintf("the field [%s] shouldn't be zero", name))
	}

	return nil
}

func middle_regex(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {

	r, e := regexp.Compile(expected.(string))
	if e != nil {
		return errors.New("0", e)
	}

	if !r.MatchString(fmt.Sprintf("%+v", value)) {
		return errors.New("0", fmt.Sprintf("the field [%s] has invalid data [%+v]", name, value))
	}

	return nil
}
