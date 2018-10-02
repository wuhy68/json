package validator

import (
	"fmt"
	"reflect"
	"strings"
)

func NewValidatorHandler(validator *Validator) *ValidatorContext {
	return &ValidatorContext{
		validator: validator,
		values:    make(map[string]*Data),
	}
}
func (v *ValidatorContext) handleValidation(obj interface{}) []error {
	errs := make([]error, 0)

	// load id's
	v.load(obj, &errs)

	// execute
	v.do(obj, &errs)

	return errs
}

func (v *ValidatorContext) do(obj interface{}, errs *[]error) error {
	types := reflect.TypeOf(obj)
	value := reflect.ValueOf(obj)

	if !value.CanInterface() {
		return nil
	}

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

			if err := v.doValidate(nextValue, nextType, errs); err != nil {

				if !v.validator.validateAll {
					return err
				}
			}
			if err := v.do(nextValue.Interface(), errs); err != nil {
				if !v.validator.validateAll {
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

			if err := v.do(nextValue.Interface(), errs); err != nil {
				if !v.validator.validateAll {
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

			if err := v.do(key.Interface(), errs); err != nil {
				if !v.validator.validateAll {
					return err
				}
			}
			if err := v.do(nextValue.Interface(), errs); err != nil {
				if !v.validator.validateAll {
					return err
				}
			}
		}

	default:
		// do nothing ...
	}
	return nil
}

func (v *ValidatorContext) load(obj interface{}, errs *[]error) error {
	types := reflect.TypeOf(obj)
	value := reflect.ValueOf(obj)

	if !value.CanInterface() {
		return nil
	}

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

			value, exists := nextType.Tag.Lookup(v.validator.tag)
			if !exists {
				continue
			}

			if strings.Contains(value, "id=") {
				split := strings.Split(value, ",")
				for _, item := range split {
					tag := strings.Split(item, "=")
					if strings.TrimSpace(tag[0]) == "id" {
						v.values[tag[1]] = &Data{
							Value: nextValue,
							Type:  nextType,
						}
					}
				}
			}

			if err := v.load(nextValue.Interface(), errs); err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)

			if !nextValue.CanInterface() {
				continue
			}

			if err := v.load(nextValue.Interface(), errs); err != nil {
				return err
			}
		}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			nextValue := value.MapIndex(key)

			if !nextValue.CanInterface() {
				continue
			}

			if err := v.do(key.Interface(), errs); err != nil {
				return err
			}
			if err := v.load(nextValue.Interface(), errs); err != nil {
				return err
			}
		}

	default:
		// do nothing ...
	}
	return nil
}

func (v *ValidatorContext) doValidate(value reflect.Value, typ reflect.StructField, errs *[]error) error {

	tag, exists := typ.Tag.Lookup(v.validator.tag)
	if !exists {
		return nil
	}

	validations := strings.Split(tag, ",")

	return v.executeHandlers(value, typ, validations, errs)
}

func (v *ValidatorContext) executeHandlers(value reflect.Value, typ reflect.StructField, validations []string, errs *[]error) error {
	var err error
	var itErrs []error
	var replacedErrors = make(map[error]bool)

	for _, validation := range validations {
		var name string

		options := strings.SplitN(validation, "=", 2)
		tag := strings.TrimSpace(options[0])

		if _, ok := v.validator.activeHandlers[tag]; !ok {
			err := fmt.Errorf("invalid tag [%s]", tag)
			*errs = append(*errs, err)

			if !v.validator.validateAll {
				return err
			}
		}

		var expected interface{}
		if len(options) > 1 {
			expected = strings.TrimSpace(options[1])
		}

		jsonName, exists := typ.Tag.Lookup("json")
		if exists {
			split := strings.SplitN(jsonName, ",", 2)
			name = split[0]
		} else {
			name = typ.Name
		}

		validationData := ValidationData{
			Name:     name,
			Value:    value,
			Expected: expected,
			Errors:   &itErrs,
			ErrorsReplaced: replacedErrors,
		}
		// execute validations
		if _, ok := v.validator.handlersBefore[tag]; ok {
			if rtnErrs := v.validator.handlersBefore[tag](v, &validationData); rtnErrs != nil && len(rtnErrs) > 0 {

				// skip validation
				if rtnErrs[0] == ErrorSkipValidation {
					return nil
				}
				itErrs = append(itErrs, rtnErrs...)
				err = rtnErrs[0]
			}
		}

		if _, ok := v.validator.handlersMiddle[tag]; ok {
			if rtnErrs := v.validator.handlersMiddle[tag](v, &validationData); rtnErrs != nil && len(rtnErrs) > 0 {
				itErrs = append(itErrs, rtnErrs...)
				err = rtnErrs[0]
			}
		}

		if _, ok := v.validator.handlersAfter[tag]; ok {
			if rtnErrs := v.validator.handlersAfter[tag](v, &validationData); rtnErrs != nil && len(rtnErrs) > 0 {
				itErrs = append(itErrs, rtnErrs...)
				err = rtnErrs[0]
			}
		}
	}

	*errs = append(*errs, itErrs...)

	return err
}
