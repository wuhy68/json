package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/joaosoft/errors"
)

func handleValidation(obj interface{}) error {
	return do(obj)
}

func do(obj interface{}) error {
	types := reflect.TypeOf(obj)
	value := reflect.ValueOf(obj)

	if value.Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem()

		if value.IsValid() {
			types = value.Type()
		} else {
			return nil
		}
	}

	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < types.NumField(); i++ {
			nextValue := value.Field(i)
			nextType := types.Field(i)

			if err := doValidate(nextValue, nextType); err != nil {
				return err
			}
			if err := do(nextValue.Interface()); err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)
			if err := do(nextValue.Interface()); err != nil {
				return err
			}
		}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			nextValue := value.MapIndex(key)
			if err := do(key.Interface()); err != nil {
				return err
			}
			if err := do(nextValue.Interface()); err != nil {
				return err
			}
		}

	default:
		// do nothing...
	}
	return nil
}

func doValidate(value reflect.Value, typ reflect.StructField) error {
	switch value.Kind() {
	case reflect.Struct:

	case reflect.Array, reflect.Slice:
		// max, min, equal

	case reflect.Map:
		// max, min, equal
	}

	tag, exists := typ.Tag.Lookup(validator.tag)
	if !exists {
		return nil
	}

	validations := strings.Split(tag, ",")
	for _, validation := range validations {
		options := strings.Split(validation, "=")
		tag := strings.TrimSpace(options[0])
		expected := strings.TrimSpace(options[1])

		if _, ok := validator.handlers[tag]; ok {
			if err := validator.handlers[tag](typ.Name, value, expected); err != nil {
				return err
			}
		} else {
			return errors.New(fmt.Sprintf("invalid tag [%s]", tag))
		}
	}

	return nil
}
