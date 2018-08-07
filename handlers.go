package validator

import (
	"fmt"

	"reflect"

	"strconv"

	"github.com/joaosoft/errors"
)

func NewDefaultHandlers() map[string]TagHandler {
	return map[string]TagHandler{"value": value, "size": size, "min": min, "max": max}
}

func value(name string, value reflect.Value, expected interface{}) error {
	if fmt.Sprintf("%+v", value) != fmt.Sprintf("%+v", expected) {
		return errors.New(fmt.Sprintf("the value [%+v] is diferent of the expected [%+v] on field [%s]", value, expected, name))
	}

	return nil
}

func size(name string, value reflect.Value, expected interface{}) error {
	size, err := strconv.Atoi(expected.(string))
	if err != nil {
		return errors.New(fmt.Sprintf("the size [%s] is invalid on field [%s]", expected, value))
	}

	if value.Len() != size {
		return errors.New(fmt.Sprintf("the value [%+v] is diferent of the expected [%+v] on field [%s]", value, expected, name))
	}

	return nil
}

func min(name string, value reflect.Value, expected interface{}) error {
	min, err := strconv.Atoi(expected.(string))
	if err != nil {
		return errors.New(fmt.Sprintf("the size [%s] is invalid on field [%s]", expected, value))
	}

	switch value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		if value.Len() < min {
			return errors.New(fmt.Sprintf("the length [%+v] is lower then the expected [%+v] on field [%s]", value, expected, name))
		}
	default:
		if value.Int() < int64(min) {
			return errors.New(fmt.Sprintf("the value [%+v] is lower then the expected [%+v] on field [%s]", value, expected, name))
		}
	}

	return nil
}

func max(name string, value reflect.Value, expected interface{}) error {
	max, err := strconv.Atoi(expected.(string))
	if err != nil {
		return errors.New(fmt.Sprintf("the size [%s] is invalid on field [%s]", expected, value))
	}

	switch value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		if value.Len() > max {
			return errors.New(fmt.Sprintf("the length [%+v] is bigger then the expected [%+v] on field [%s]", value, expected, name))
		}
	default:
		if value.Int() < int64(max) {
			return errors.New(fmt.Sprintf("the value [%+v] is bigger then the expected [%+v] on field [%s]", value, expected, name))
		}
	}

	return nil
}
