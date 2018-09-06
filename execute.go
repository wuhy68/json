package validator

import (
	"fmt"
	"reflect"
	"strings"
)

func handleValidation(obj interface{}) []error {
	errs := make([]error, 0)

	do(obj, &errs)

	return errs
}

func do(obj interface{}, errs *[]error) error {
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

			if !nextValue.CanInterface() {
				continue
			}

			if err := doValidate(nextValue.Interface(), nextType, errs); err != nil {

				if !validatorInstance.validateAll {
					return err
				}
			}
			if err := do(nextValue.Interface(), errs); err != nil {
				if !validatorInstance.validateAll {
					return err
				}
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)

			if !nextValue.CanInterface() {
				continue
			}

			if err := do(nextValue.Interface(), errs); err != nil {
				if !validatorInstance.validateAll {
					return err
				}
			}
		}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			nextValue := value.MapIndex(key)

			if !nextValue.CanInterface() {
				continue
			}

			if err := do(key.Interface(), errs); err != nil {
				if !validatorInstance.validateAll {
					return err
				}
			}
			if err := do(nextValue.Interface(), errs); err != nil {
				if !validatorInstance.validateAll {
					return err
				}
			}
		}

	default:
		// do nothing ...
	}
	return nil
}

func doValidate(value interface{}, typ reflect.StructField, errs *[]error) error {

	tag, exists := typ.Tag.Lookup(validatorInstance.tag)
	if !exists {
		return nil
	}

	validations := strings.Split(tag, ",")

	return executeHandlers(reflect.ValueOf(value), typ, validations, errs)
}

func executeHandlers(value reflect.Value, typ reflect.StructField, validations []string, errs *[]error) error {
	var err error
	var itErrs []error

	for _, validation := range validations {
		var name string

		options := strings.Split(validation, "=")
		tag := strings.TrimSpace(options[0])

		if _, ok := validatorInstance.activeHandlers[tag]; !ok {
			err := fmt.Errorf("invalid tag [%s]", tag)
			*errs = append(*errs, err)

			if !validatorInstance.validateAll {
				return err
			}
		}

		var expected string
		if len(options) > 1 {
			expected = strings.TrimSpace(options[1])
		}

		jsonName, exists := typ.Tag.Lookup("json")
		if exists {
			name = jsonName
		} else {
			name = typ.Name
		}

		if _, ok := validatorInstance.handlersBefore[tag]; ok {
			if rtnErrs := validatorInstance.handlersBefore[tag](name, value, expected); len(rtnErrs) > 0 {
				itErrs = append(itErrs, rtnErrs...)
				err = rtnErrs[0]
			}
		}

		if _, ok := validatorInstance.handlersMiddle[tag]; ok {
			if rtnErrs := validatorInstance.handlersMiddle[tag](name, value, expected, &itErrs); len(rtnErrs) > 0 {
				itErrs = append(itErrs, rtnErrs...)
				err = rtnErrs[0]
			}
		}

		if _, ok := validatorInstance.handlersAfter[tag]; ok {
			if rtnErrs := validatorInstance.handlersAfter[tag](name, value, expected, &itErrs); len(rtnErrs) > 0 {
				itErrs = append(itErrs, rtnErrs...)
				err = rtnErrs[0]
			}
		}
	}

	*errs = append(*errs, itErrs...)

	return err
}
