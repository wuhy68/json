package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/joaosoft/errors"
)

func handleValidation(obj interface{}) *errors.ListErr {
	var errs errors.ListErr

	do(obj, &errs)

	return &errs
}

func do(obj interface{}, errs *errors.ListErr) error {
	types := reflect.TypeOf(obj)
	value := reflect.ValueOf(obj)

	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()

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

			if nextValue.Kind() == reflect.Ptr && !nextValue.IsNil() {
				nextValue = nextValue.Elem()
			}

			if err := doValidate(nextValue.Interface(), nextType, errs); err != nil {

				if !validator.validateAll {
					return err
				}
			}
			if err := do(nextValue.Interface(), errs); err != nil {
				if !validator.validateAll {
					return err
				}
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)
			if err := do(nextValue.Interface(), errs); err != nil {
				if !validator.validateAll {
					return err
				}
			}
		}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			nextValue := value.MapIndex(key)
			if err := do(key.Interface(), errs); err != nil {
				if !validator.validateAll {
					return err
				}
			}
			if err := do(nextValue.Interface(), errs); err != nil {
				if !validator.validateAll {
					return err
				}
			}
		}

	default:
		// do nothing ...
	}
	return nil
}

func doValidate(value interface{}, typ reflect.StructField, errs *errors.ListErr) error {

	tag, exists := typ.Tag.Lookup(validator.tag)
	if !exists {
		return nil
	}

	validations := strings.Split(tag, ",")

	return executeHandlers(reflect.ValueOf(value), typ, validations, errs)
}

func executeHandlers(value reflect.Value, typ reflect.StructField, validations []string, errs *errors.ListErr) error {
	var err error
	var itErrs errors.ListErr

	for _, validation := range validations {
		options := strings.Split(validation, "=")

		tag := strings.TrimSpace(options[0])

		if _, ok := validator.activeHandlers[tag]; !ok {
			err := errors.New("0", fmt.Sprintf("invalid tag [%s]", tag))
			*errs = append(*errs, err)

			if !validator.validateAll {
				return err
			}
		}

		var expected string
		if len(options) > 1 {
			expected = strings.TrimSpace(options[1])
		}

		if _, ok := validator.handlersPre[tag]; ok {
			if rtnErrs := validator.handlersPre[tag](typ.Name, value, expected); !rtnErrs.IsEmpty() {
				itErrs = append(itErrs, rtnErrs...)
				err = rtnErrs[0]
			}
		}

		if _, ok := validator.handlersMiddle[tag]; ok {
			if rtnErrs := validator.handlersMiddle[tag](typ.Name, value, expected, &itErrs); !rtnErrs.IsEmpty() {
				itErrs = append(itErrs, rtnErrs...)
				err = rtnErrs[0]
			}
		}

		if _, ok := validator.handlersPos[tag]; ok {
			if rtnErrs := validator.handlersPos[tag](typ.Name, value, expected, &itErrs); !rtnErrs.IsEmpty() {
				itErrs = append(itErrs, rtnErrs...)
				err = rtnErrs[0]
			}
		}
	}

	*errs = append(*errs, itErrs...)

	return err
}
