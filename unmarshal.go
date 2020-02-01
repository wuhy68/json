package json

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var typeUnmarshal = reflect.TypeOf((*iunmarshal)(nil)).Elem()

type iunmarshal interface {
	UnmarshalJSON(bytes []byte) error
}

type unmarshal struct {
	bytes  []byte
	object interface{}
	tags   []string
}

func newUnmarshal(bytes []byte, object interface{}, tags ...string) *unmarshal {
	return &unmarshal{bytes: bytes, object: object, tags: tags}
}

func (u *unmarshal) execute() error {
	value := reflect.ValueOf(u.object)

	if value.Kind() != reflect.Ptr || value.IsNil() {
		return ErrorInvalidPointer
	}

	return u.do(value, u.bytes, 0)
}

func (u *unmarshal) do(object reflect.Value, byts []byte, iteration int) error {
	iteration++
	var err error
	exit := false

	for len(byts) > 0 && !exit {
		switch string(byts[0:1]) {
		case comma:
			err = u.handle(object, byts[1:], iteration)
			if err != nil {
				return err
			}
			exit = true

		case jsonStart, arrayStart:
			err = u.handle(object, byts, iteration)
			if err != nil {
				return err
			}
			exit = true

		default:
			byts = byts[1:]
			continue
		}
	}
	return nil
}

func (u *unmarshal) handle(object reflect.Value, byts []byte, iteration int) error {
	var fieldName, fieldValue, nextValue []byte
	var sliceValues [][]byte
	var err error
	originByts := byts

	switch string(byts[0]) {
	case arrayStart:
		fieldValue, _, nextValue, err = u.getJsonValue(byts)
		if err != nil {
			return err
		}

		sliceValues, _, err = u.getJsonSliceValues(fieldValue)
		if err != nil {
			return err
		}
	case jsonStart:
		byts = byts[1 : len(byts)-1]
		fallthrough
	default:
		fieldName, byts, err = u.getJsonName(byts)
		if err != nil {
			return err
		}

		switch string(byts[0]) {
		case arrayStart:
			fieldValue, _, nextValue, err = u.getJsonValue(byts)
			if err != nil {
				return err
			}

			sliceValues, _, err = u.getJsonSliceValues(fieldValue)
			if err != nil {
				return err
			}
		default:
			fieldValue, _, nextValue, err = u.getJsonValue(byts)
			if err != nil {
				return err
			}
		}
	}

	exists, field, err := u.getField(object, string(fieldName))
	if err != nil {
		return err
	}

	if exists {
		isPrimitive, kind, err := u.isPrimitive(field)
		if err != nil {
			return err
		}

		if isPrimitive {
			if err = u.setField(field, string(fieldValue)); err != nil {
				return err
			}
		} else {
			newUnmarshalField, handled, err := u.handleUnmarshalJSON(field, fieldValue)
			if err != nil {
				return err
			}

			if handled {
				field.Set(newUnmarshalField)
				goto next
			}

			if !bytes.HasPrefix(fieldValue, []byte(null)) {
				if field.Kind() == reflect.Ptr && field.IsNil() {
					field.Set(reflectAlloc(field.Type()))
				}

				if field.Kind() == reflect.Ptr {
					field = field.Elem()
				}
			}

			// is a slice
			switch kind {
			case reflect.Array, reflect.Slice:
				lenField := len(sliceValues)
				field.Set(reflect.MakeSlice(field.Type(), 0, lenField))

				for _, item := range sliceValues {

					newValue := reflectAlloc(field.Type().Elem())

					isPrimitive, kind, err = u.isPrimitive(newValue)
					if err != nil {
						return err
					}

					newUnmarshalValue, handled, err := u.handleUnmarshalJSON(newValue, item)
					if err != nil {
						return err
					}

					if handled {
						newValue.Set(newUnmarshalValue)
						continue
					}

					if isPrimitive {
						if err = u.setField(newValue, string(item)); err != nil {
							return err
						}
					} else {
						if err = u.do(newValue, item, iteration); err != nil {
							return err
						}
					}

					field.Set(reflect.Append(field, newValue))
				}

			case reflect.Map:
				field.Set(reflect.MakeMap(field.Type()))

				if iteration == 1 {
					byts = originByts
					nextValue = []byte{}
				}

				mapValues, keysType, valuesType, err := u.getJsonMapValues(byts)
				if err != nil {
					return err
				}

				for key, value := range mapValues {
					// key
					newKeyType := reflect.New(field.Type().Key()).Elem().Type()
					if newKeyType.Kind() == reflect.Interface {
						newKeyType = keysType[key]
					}
					newKey := reflectAlloc(newKeyType)

					isPrimitive, kind, err = u.isPrimitive(newKey)
					if err != nil {
						return err
					}

					newUnmarshalField, handled, err := u.handleUnmarshalJSON(newKey, []byte(key))
					if err != nil {
						return err
					}

					if handled {
						newKey.Set(newUnmarshalField)
						continue
					}
					if isPrimitive {
						if err = u.setField(newKey, key); err != nil {
							return err
						}
					} else {
						if err = u.do(newKey, []byte(key), iteration); err != nil {
							return err
						}
					}

					// value
					newValueType := field.Type().Elem()
					if newValueType.Kind() == reflect.Interface {
						newValueType = valuesType[key]
					}
					newValue := reflectAlloc(newValueType)

					isPrimitive, kind, err = u.isPrimitive(newValue)
					if err != nil {
						return err
					}

					newValue, handled, err = u.handleUnmarshalJSON(newValue, value)
					if err != nil {
						return err
					}

					if !handled {
						if isPrimitive {
							if err = u.setField(newValue, string(value)); err != nil {
								return err
							}
						} else {
							if err = u.do(newValue, []byte(value), iteration); err != nil {
								return err
							}
						}
					}

					field.SetMapIndex(newKey, newValue)
				}
			default:
				if err = u.do(field, fieldValue, iteration); err != nil {
					return err
				}
			}
		}
	}

next:
	// next
	if err = u.do(object, nextValue, iteration); err != nil {
		return err
	}

	return nil
}

func (u *unmarshal) handleUnmarshalJSON(object reflect.Value, byts []byte) (reflect.Value, bool, error) {
	val, ok := object.Interface().(iunmarshal)
	if ok {
		if err := val.UnmarshalJSON([]byte(`"2020-02-01T09:59:32.634417Z"`)); err != nil {
			return object, true, err
		}
		return reflect.ValueOf(val), true, nil
	}

	return object, ok, nil
}

func (u *unmarshal) getField(object reflect.Value, name string) (bool, reflect.Value, error) {
	hasImplementation := object.Type().Implements(typeUnmarshal)
	isSlice := object.Kind() == reflect.Slice && !hasImplementation
	isMap := object.Kind() == reflect.Map && !hasImplementation
	isMapOfSlices := isMap && object.Type().Elem().Kind() == reflect.Slice

	if isMapOfSlices && object.IsNil() {
		object.Set(reflectAlloc(object.Type().Elem().Elem()))
	} else if (isSlice || isMap) && object.IsNil() {
		object.Set(reflectAlloc(object.Type().Elem()))
	} else if object.Kind() == reflect.Ptr && object.IsNil() {
		object.Set(reflectAlloc(object.Type().Elem()))
	}

	types := reflect.TypeOf(object.Interface())
	innerObject := object

	if !innerObject.CanInterface() {
		return false, object, nil
	}

	if innerObject.Kind() == reflect.Ptr && !innerObject.IsNil() {
		innerObject = innerObject.Elem()

		if innerObject.IsValid() {
			types = innerObject.Type()
		} else {
			return false, object, nil
		}
	}

	switch innerObject.Kind() {
	case reflect.Struct:
		for i := 0; i < types.NumField(); i++ {
			nextValue := innerObject.Field(i)
			object = nextValue
			nextType := types.Field(i)

			if nextValue.Kind() == reflect.Ptr {
				if !nextValue.IsNil() {
					nextValue = nextValue.Elem()
				} else {
					nextValue.Set(reflectAlloc(nextValue.Type()))
				}
			}

			if !nextValue.CanInterface() {
				continue
			}

			exists, tag, err := u.loadTag(nextType)
			if err != nil {
				return false, object, err
			}

			if exists && tag == name {
				return true, object, nil
			}
		}

	case reflect.Array, reflect.Slice:
		return true, object, nil

	case reflect.Map:
		return true, object, nil

	default:
	}

	return false, object, nil
}

func (u *unmarshal) setField(object reflect.Value, value string) error {
	if object.Kind() == reflect.Ptr {
		if !object.IsNil() {
			object = object.Elem()
		} else {
			object = reflectAlloc(object.Type()).Elem()
		}
	}

	switch object.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tmpValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		object.SetInt(tmpValue)
	case reflect.Float32, reflect.Float64:
		v, _ := strconv.ParseFloat(value, 64)
		object.SetFloat(v)
	case reflect.Bool:
		v, _ := strconv.ParseBool(value)
		object.SetBool(v)
	case reflect.String:
		object.SetString(u.decodeString(value))
	}

	return nil
}

func (u *unmarshal) isPrimitive(object reflect.Value) (bool, reflect.Kind, error) {

	if object.Kind() == reflect.Ptr {
		if !object.IsNil() {
			object = object.Elem()
		} else {
			object = reflectAlloc(object.Type()).Elem()
		}
	}

	switch object.Kind() {
	case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
		return false, object.Kind(), nil
	default:
		return true, object.Kind(), nil
	}
}

func (u *unmarshal) getJsonName(byts []byte) ([]byte, []byte, error) {

	start := bytes.Index(byts, []byte(stringStartEnd))
	end := bytes.Index(byts[start+1:], []byte(stringStartEnd))

	startNext := bytes.Index(byts[start+end:], []byte(is))

	fieldName := byts[start+1 : start+end+1]

	return fieldName, byts[start+end+startNext+1:], nil
}

func (u *unmarshal) getJsonValue(byts []byte) (value []byte, valueType reflect.Type, nextValue []byte, err error) {
	startInit := false
	startEnd := false
	numJsonStart := 0
	start := 0
	end := 0
	open := false
	exit := false

	for i, field := range byts {
		switch string(field) {

		case comma:
			if !open && numJsonStart == 0 {
				end = i
				exit = true
			}

		case stringStartEnd:
			valueType = TypeString
			if i > 0 {
				if byts[i-1] == []byte(`\`)[0] {
					continue
				}
			}

			if !open {
				if !startInit {
					startInit = true
					start = i + 1
				}
				numJsonStart++
			} else {
				numJsonStart--
				startEnd = true
				end = i
			}

			open = !open

		case jsonStart:
			if !startInit {
				startInit = true
				start = i
			}

			numJsonStart++

		case jsonEnd:
			startEnd = true
			end = i + 1
			numJsonStart--

		case arrayStart:
			if !startInit {
				startInit = true
				start = i + 1
			}

			numJsonStart++

		case arrayEnd:
			startEnd = true
			end = i
			numJsonStart--

		default:
			end = i + 1
			continue
		}

		if numJsonStart <= 0 || exit {
			if !startEnd {
				end = i
			}
			break
		}
	}

	nextEnd := end
	if len(byts) > end {
		if byts[end] == []byte(stringStartEnd)[0] {
			nextEnd = end + 1
		}
	}

	item := bytes.TrimSpace(byts[start:end])

	switch valueType {
	case TypeString:
		if _, err = time.Parse(time.RFC3339Nano, string(item)); err == nil {
			valueType = TypeTimestamp
		}
	default:
		switch string(item) {
		case booleanTrue, booleanFalse:
			valueType = TypeBoolean
		default:
			if _, err = strconv.ParseInt(string(item), 10, 64); err == nil {
				valueType = TypeFloat64
			} else if _, err = strconv.ParseFloat(string(item), 64); err == nil {
				valueType = TypeFloat64
			} else {
				valueType = TypeString
			}
		}
	}

	return item, valueType, byts[nextEnd:], nil
}

func (u *unmarshal) getJsonSliceValues(byts []byte) (values [][]byte, valuesType []reflect.Type, err error) {

	var value []byte
	var valueType reflect.Type

	for len(byts) > 0 {
		value, valueType, byts, err = u.getJsonValue(byts)
		if err != nil {
			return nil, nil, err
		}
		values = append(values, value)
		valuesType = append(valuesType, valueType)

		exit := false
		for len(byts) > 0 && !exit {
			switch string(byts[0]) {
			case comma, arrayEnd, empty:
				byts = byts[1:]
				if len(byts) == 0 {
					exit = true
				}
			default:
				exit = true
			}
		}
	}

	return values, valuesType, nil
}

func (u *unmarshal) getJsonMapValues(byts []byte) (_ map[string][]byte, keysType map[string]reflect.Type, valuesType map[string]reflect.Type, err error) {

	var key []byte
	var value []byte
	var values = make(map[string][]byte)

	keysType = make(map[string]reflect.Type)
	valuesType = make(map[string]reflect.Type)

	if string(byts[0]) == jsonStart && string(byts[len(byts)-1]) == jsonEnd {
		byts = byts[1 : len(byts)-1]
	}

	var itemType reflect.Type
	for len(byts) > 0 {
		key, itemType, byts, err = u.getJsonValue(byts)
		if err != nil {
			return nil, nil, nil, err
		}

		strKey := string(key)
		keysType[strKey] = itemType

		exit := false
		for !exit && len(byts) > 0 {
			switch string(byts[0]) {
			case is:
				byts = byts[1:]
				exit = true
			default:
				byts = byts[1:]
			}
		}

		value, itemType, byts, err = u.getJsonValue(byts)
		if err != nil {
			return nil, nil, nil, err
		}

		values[string(key)] = value
		valuesType[strKey] = itemType

		exit = false
		for !exit && len(byts) > 0 {
			switch string(byts[0]) {
			case comma:
				byts = byts[1:]
				exit = true
			default:
				byts = byts[1:]
			}
		}
	}

	return values, keysType, valuesType, nil
}

func (u *unmarshal) loadTag(typ reflect.StructField) (exists bool, tag string, err error) {
	for _, searchTag := range u.tags {
		tag, exists = typ.Tag.Lookup(searchTag)

		if exists && tag == "-" {
			tag = ""
			exists = false
			break
		}

		if exists {
			break
		}
	}

	return exists, tag, err
}

func (u *unmarshal) decodeString(str string) string {
	return strings.Replace(str, stringStartEndScaped, stringStartEnd, -1)
}
