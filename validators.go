package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"encoding/json"
	"errors"

	"github.com/satori/go.uuid"
)

func (v *Validator) validate_value(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	rtnErrs := make([]error, 0)

	if fmt.Sprintf("%+v", value) == "" || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return rtnErrs
	}

	if fmt.Sprintf("%+v", value) != fmt.Sprintf("%+v", expected) {
		err := fmt.Errorf("the value [%+v] is different of the expected [%+v] on field [%s]", value, expected, name)
		rtnErrs = append(rtnErrs, err)
	}

	return rtnErrs
}

func (v *Validator) validate_sanitize(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	rtnErrs := make([]error, 0)

	if fmt.Sprintf("%+v", value) == "" || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return rtnErrs
	}

	val := fmt.Sprintf("%+v", value)
	split := strings.Split(expected.(string), ";")
	invalid := make([]string, 0)

	// validate expected
	for _, str := range split {
		if strings.Contains(val, str) {
			invalid = append(invalid, str)
		}
	}

	// validate global
	for _, str := range v.sanitize {
		if strings.Contains(val, str) {
			invalid = append(invalid, str)
		}
	}

	if len(invalid) > 0 {
		err := fmt.Errorf("the value [%+v] is has invalid characters [%+v] on field [%s]", value, strings.Join(invalid, ","), name)
		rtnErrs = append(rtnErrs, err)
	}

	return rtnErrs
}

func (v *Validator) validate_not(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	rtnErrs := make([]error, 0)

	if fmt.Sprintf("%+v", value) == "" || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return rtnErrs
	}

	if fmt.Sprintf("%+v", value) == fmt.Sprintf("%+v", expected) {
		err := fmt.Errorf("the value [%+v] should be different of the [%+v] on field [%s]", value, expected, name)
		rtnErrs = append(rtnErrs, err)
	}

	return rtnErrs
}

func (v *Validator) validate_options(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	rtnErrs := make([]error, 0)

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

			if !nextValue.CanInterface() {
				continue
			}

			_, ok := optionsVal[fmt.Sprintf("%+v", nextValue.Interface())]
			if !ok {
				invalidValue = nextValue.Interface()
				err := fmt.Errorf("the value [%+v] is different of the expected options [%+v] on field [%s]", invalidValue, expected, name)
				rtnErrs = append(rtnErrs, err)
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

			if !nextValue.CanInterface() {
				continue
			}

			val, ok := optionsMap[fmt.Sprintf("%+v", key.Interface())]
			if !ok || fmt.Sprintf("%+v", nextValue.Interface()) != fmt.Sprintf("%+v", val) {
				invalidValue = fmt.Sprintf("%+v:%+v", key.Interface(), nextValue.Interface())
				err := fmt.Errorf("the value [%+v] is different of the expected options [%+v] on field [%s]", nextValue.Interface(), expected, name)
				rtnErrs = append(rtnErrs, err)
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
			err := fmt.Errorf("the value [%+v] is different of the expected options [%+v] on field [%s]", invalidValue, expected, name)
			rtnErrs = append(rtnErrs, err)
		}
	}

	return rtnErrs
}

func (v *Validator) validate_size(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	rtnErrs := make([]error, 0)
	size, e := strconv.Atoi(expected.(string))
	if e != nil {
		err := fmt.Errorf("the size [%s] is invalid on field [%s]", expected, value)
		rtnErrs = append(rtnErrs, err)
		return rtnErrs
	}

	var valueSize int64

	switch value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		valueSize = int64(value.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valueSize = int64(len(strings.TrimSpace(strconv.Itoa(int(value.Int())))))
	case reflect.Float32, reflect.Float64:
		valueSize = int64(len(strings.TrimSpace(strconv.FormatFloat(value.Float(), 'g', 1, 64))))
	case reflect.String:
		valueSize = int64(len(strings.TrimSpace(value.String())))
	case reflect.Bool:
		valueSize = int64(len(strings.TrimSpace(strconv.FormatBool(value.Bool()))))
	default:
		if value.Kind() == reflect.Ptr && value.IsNil() {
			break
		}
		valueSize = int64(len(strings.TrimSpace(value.String())))
	}

	if valueSize != int64(size) {
		err := fmt.Errorf("the length [%+v] is lower then the expected [%+v] on field [%s]", valueSize, expected, name)
		rtnErrs = append(rtnErrs, err)
	}

	return rtnErrs
}

func (v *Validator) validate_min(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	rtnErrs := make([]error, 0)
	min, e := strconv.Atoi(expected.(string))
	if e != nil {
		err := fmt.Errorf("the size [%s] is invalid on field [%s]", expected, value)
		rtnErrs = append(rtnErrs, err)
		return rtnErrs
	}

	var valueSize int64

	switch value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		valueSize = int64(value.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valueSize = value.Int()
	case reflect.Float32, reflect.Float64:
		valueSize = int64(value.Float())
	case reflect.String:
		valueSize = int64(len(strings.TrimSpace(value.String())))
	case reflect.Bool:
		valueSize = int64(len(strings.TrimSpace(strconv.FormatBool(value.Bool()))))
	default:
		if value.Kind() == reflect.Ptr && value.IsNil() {
			break
		}
		valueSize = int64(len(strings.TrimSpace(value.String())))
	}

	if valueSize < int64(min) {
		err := fmt.Errorf("the length [%+v] is lower then the expected [%+v] on field [%s]", valueSize, expected, name)
		rtnErrs = append(rtnErrs, err)
	}

	return rtnErrs
}

func (v *Validator) validate_max(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	rtnErrs := make([]error, 0)
	max, e := strconv.Atoi(expected.(string))
	if e != nil {
		err := fmt.Errorf("the size [%s] is invalid on field [%s]", expected, value)
		rtnErrs = append(rtnErrs, err)
		return rtnErrs
	}

	var valueSize int64

	switch value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		valueSize = int64(value.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valueSize = value.Int()
	case reflect.Float32, reflect.Float64:
		valueSize = int64(value.Float())
	case reflect.String:
		valueSize = int64(len(strings.TrimSpace(value.String())))
	case reflect.Bool:
		valueSize = int64(len(strings.TrimSpace(strconv.FormatBool(value.Bool()))))
	default:
		if value.Kind() == reflect.Ptr && value.IsNil() {
			break
		}
		valueSize = int64(len(strings.TrimSpace(value.String())))
	}

	if valueSize > int64(max) {
		err := fmt.Errorf("the length [%+v] is bigger then the expected [%+v] on field [%s]", valueSize, expected, name)
		rtnErrs = append(rtnErrs, err)
	}

	return rtnErrs
}

func (v *Validator) validate_nonzero(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	rtnErrs := make([]error, 0)
	var valueSize int64
	var val string

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

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = strings.TrimSpace(strconv.Itoa(int(value.Int())))
		valueSize = int64(len(val))
	case reflect.Float32, reflect.Float64:
		val = strings.TrimSpace(strconv.FormatFloat(value.Float(), 'g', 1, 64))
		valueSize = int64(len(val))
	case reflect.String:
		valueSize = int64(len(strings.TrimSpace(value.String())))
	case reflect.Bool:
		valueSize = int64(len(strings.TrimSpace(strconv.FormatBool(value.Bool()))))
	case reflect.Struct:
		if value.Interface() != reflect.Zero(value.Type()).Interface() {
			valueSize = 1
		}
	default:
		if value.Kind() == reflect.Ptr && value.IsNil() {
			break
		}
		valueSize = int64(len(strings.TrimSpace(value.String())))
	}

	if valueSize == 0 || (val == "0") {
		err := fmt.Errorf("the value shouldn't be zero on field [%s]", name)
		rtnErrs = append(rtnErrs, err)
	}

	return rtnErrs
}

func (v *Validator) validate_regex(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {

	rtnErrs := make([]error, 0)

	if value.Kind() == reflect.Ptr && value.IsNil() {
		return rtnErrs
	}

	if fmt.Sprintf("%+v", value) == "" || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return rtnErrs
	}

	r, err := regexp.Compile(expected.(string))
	if err != nil {
		rtnErrs = append(rtnErrs, err)
		return rtnErrs
	}

	if len(fmt.Sprintf("%+v", value)) > 0 {
		if !r.MatchString(fmt.Sprintf("%+v", value)) {
			err := fmt.Errorf("invalid data [%s] on field [%+v] ", value, name)
			rtnErrs = append(rtnErrs, err)
		}
	}

	return rtnErrs
}

func (v *Validator) validate_special(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {

	rtnErrs := make([]error, 0)

	if fmt.Sprintf("%+v", value) == "" || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return rtnErrs
	}

	switch expected {
	case TagForDateDefault:
		expected = RegexForDateDefault
	case TagForDateDDMMYYYY:
		expected = RegexForDateDDMMYYYY
	case TagForDateYYYYMMDD:
		expected = RegexForDateYYYYMMDD
	case TagForTimeDefault:
		expected = RegexForTimeDefault
	case TagForTimeHHMMSS:
		expected = RegexForTimeHHMMSS
	case TagForURL:
		expected = RegexForURL
	default:
		err := fmt.Errorf("invalid special [%s] on field [%+v] ", expected, name)
		rtnErrs = append(rtnErrs, err)
		return rtnErrs
	}

	return v.validate_regex(context, name, value, expected, errs)
}

func (v *Validator) validate_callback(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {

	if callback, ok := v.callbacks[expected.(string)]; ok {
		return callback(context, name, value, expected, errs)
	}

	return make([]error, 0)
}

func (v *Validator) validate_error(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	rtnErrs := make([]error, 0)
	added := make(map[string]bool)
	for i, _ := range *errs {
		if v.errorCodeHandler != nil {
			if matched, err := regexp.MatchString(RegexForErrorTag, expected.(string)); err != nil {
				rtnErrs = append(rtnErrs, err)
			} else {

				if matched {
					replacer := strings.NewReplacer("{", "", "}", "")
					expected := replacer.Replace(expected.(string))

					split := strings.SplitN(expected, ":", 2)
					if len(split) == 0 {
						rtnErrs = append(rtnErrs, fmt.Errorf("invalid tag error defined %s", expected))
						continue
					}

					if _, ok := added[split[0]]; !ok {
						var arguments []interface{}
						if len(split) == 2 {
							splitArgs := strings.Split(split[1], ";")
							for _, arg := range splitArgs {
								arguments = append(arguments, arg)
							}
						}
						newErr := v.errorCodeHandler(context, split[0], arguments, name, value, expected, errs)
						if newErr != nil {
							(*errs)[i] = newErr
						}

						added[split[0]] = true
					} else {
						*errs = append((*errs)[:i], (*errs)[i+1:]...)
					}
				} else {
					messageBytes, _ := json.Marshal(Error{
						Code:    fmt.Sprintf("%+v", expected),
						Message: (*errs)[i].Error(),
					})
					(*errs)[i] = errors.New(string(messageBytes))
				}
			}
		}
	}

	return rtnErrs
}

func (v *Validator) validate_match(context *ValidatorContext, name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	if expectedValue, ok := context.values[expected.(string)]; ok {
		expected = expectedValue.Value
	}

	return v.validate_value(context, name, value, expected, errs)
}

func (v *Validator) validate_id(context *ValidatorContext, name string, value reflect.Value, expected interface{}) []error {
	return nil
}

func (v *Validator) validate_if(context *ValidatorContext, name string, value reflect.Value, expected interface{}) []error {
	rtnErrs := make([]error, 0)

	str := expected.(string)
	var expressions []*Expression
	var expression *Expression
	var query string

	// read conditions
	size := len(str)

	for i := 0; i < size; i++ {
		switch str[i] {
		case '(':
			continue

		case ')':
			start := strings.Index(query, "id=")
			if start == -1 {
				return rtnErrs
			}

			end := strings.Index(query[start:], " ")
			if end == -1 {
				end = size - 1
			}

			id := query[start+3 : end]
			query = query[end+1:]

			if data, ok := context.values[id]; ok {
				var errs []error
				err := context.executeHandlers(data.Value, data.Type, strings.Split(query, " "), &errs)

				// get next operator
				var operator Operator
				if index := strings.Index(str[i+1:], "("); index > -1 {
					operator = Operator(strings.TrimSpace(str[i+1 : i+1+index]))

					str = str[i+1+index:]
					i = 0
					size = len(str)
				}

				expression = &Expression{
					Data:         data,
					Result:       err,
					NextOperator: operator,
					Expected:     query,
				}
				expressions = append(expressions, expression)
			}
			query = ""

		default:
			query = fmt.Sprintf("%s%c", query, str[i])
		}
	}

	// validate all conditions
	var condition = ""
	var prevOperator = NONE

	for _, expr := range expressions {

		if condition == "" {
			if expr.Result == nil {
				condition = "ok"
			} else {
				condition = "ko"
			}
		} else {

			switch prevOperator {
			case AND:
				if expr.Result != nil {
					condition = "ko"
				}
			case OR:
				if expr.Result == nil && condition == "ko" {
					condition = "ok"
				}
			case NONE:
				if expr.Result == nil {
					condition = "ok"
				}
			}
		}

		prevOperator = expr.NextOperator
	}

	if condition == "ko" {
		return []error{ErrorSkipValidation}
	}

	return nil
}
