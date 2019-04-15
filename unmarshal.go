package json

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
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

	return u.do(value, u.bytes)
}

func (u *unmarshal) do(object reflect.Value, byts []byte) error {
	var err error
	exit := false

	for len(byts) > 0 && !exit {
		switch string(byts[0:1]) {
		case comma:
			err = u.handle(object, byts[1:])
			if err != nil {
				return err
			}
			exit = true

		case jsonStart:
			err = u.handle(object, byts[1:])
			if err != nil {
				return err
			}
			exit = true

		case arrayStart:
			err = u.handle(object, byts[1:])
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

func (u *unmarshal) handle(object reflect.Value, byts []byte) error {
	var fieldName []byte
	var err error

	fieldName, byts, err = u.getJsonName(byts)
	if err != nil {
		return err
	}

	fieldValue, nextValue, err := u.getJsonValue(byts)
	if err != nil {
		return err
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
				sliceValues, err := u.getJsonSliceValues(fieldValue)
				if err != nil {
					return err
				}

				lenField := len(sliceValues)
				field.Set(reflect.MakeSlice(field.Type(), 0, lenField))

				for _, item := range sliceValues {

					newValue := reflectAlloc(field.Type().Elem())

					isPrimitive, kind, err = u.isPrimitive(newValue)
					if err != nil {
						return err
					}

					if isPrimitive {
						if err = u.setField(newValue, string(item)); err != nil {
							return err
						}
					} else {
						if err = u.do(newValue, item); err != nil {
							return err
						}
					}

					field.Set(reflect.Append(field, newValue))
				}
			case reflect.Map:
				mapValues, err := u.getJsonMapValues(fieldValue)
				if err != nil {
					return err
				}

				field.Set(reflect.MakeMap(field.Type()))

				for key, value := range mapValues {

					// key
					newKey := reflect.New(field.Type().Key()).Elem()
					isPrimitive, kind, err = u.isPrimitive(newKey)
					if err != nil {
						return err
					}

					if isPrimitive {
						if err = u.setField(newKey, key); err != nil {
							return err
						}
					} else {
						if err = u.do(newKey, []byte(key)); err != nil {
							return err
						}
					}

					// value
					newValue := reflectAlloc(field.Type().Elem())
					if newValue.Kind() == reflect.Ptr {
					}

					isPrimitive, kind, err = u.isPrimitive(newValue)
					if err != nil {
						return err
					}

					if isPrimitive {
						if err = u.setField(newValue, string(value)); err != nil {
							return err
						}
					} else {
						if err = u.do(newValue, []byte(value)); err != nil {
							return err
						}
					}
					field.SetMapIndex(newKey, newValue)
				}
			default:
				if err = u.do(field, fieldValue); err != nil {
					return err
				}
			}
		}
	}

	// next
	if err = u.do(object, nextValue); err != nil {
		return err
	}

	return nil
}

func (u *unmarshal) getField(object reflect.Value, name string) (bool, reflect.Value, error) {
	hasImplementation := object.Type().Implements(typeUnmarshal)
	isSlice := object.Kind() == reflect.Slice && !hasImplementation
	isMap := object.Kind() == reflect.Map && !hasImplementation
	isMapOfSlices := isMap && object.Type().Elem().Kind() == reflect.Slice

	if isMapOfSlices {
		object.Set(reflectAlloc(object.Type().Elem().Elem()))
	} else if isSlice || isMap {
		object.Set(reflectAlloc(object.Type().Elem()))
	}

	types := reflect.TypeOf(object.Interface())

	if !object.CanInterface() {
		return false, object, nil
	}

	if object.Kind() == reflect.Ptr && !object.IsNil() {
		object = object.Elem()

		if object.IsValid() {
			types = object.Type()
		} else {
			return false, object, nil
		}
	}

	switch object.Kind() {
	case reflect.Struct:
		if hasImplementation {
			if val, ok := object.Interface().(iunmarshal); ok {

				if err := val.UnmarshalJSON([]byte{}); err != nil {
					return false, object, err
				}

				return false, object, nil
			}
		}

		for i := 0; i < types.NumField(); i++ {
			nextValue := object.Field(i)
			nextType := types.Field(i)

			if nextValue.Kind() == reflect.Ptr && !nextValue.IsNil() {
				nextValue = nextValue.Elem()
			}

			if !nextValue.CanInterface() {
				continue
			}

			exists, tag, err := u.loadTag(nextValue, nextType)
			if err != nil {
				return false, object, err
			}

			if exists && tag == name {
				return true, nextValue, nil
			}
		}

	case reflect.Array, reflect.Slice:
		if object.IsNil() {
			return false, object, nil
		}

		for i := 0; i < object.Len(); i++ {
			nextValue := object.Index(i)

			if !nextValue.CanInterface() {
				continue
			}
		}

	case reflect.Map:
		if object.IsNil() {
			return false, object, nil
		}

		for _, key := range object.MapKeys() {
			nextValue := object.MapIndex(key)

			if !nextValue.CanInterface() {
				continue
			}
		}

	default:
	}
	return false, object, nil
}

func (u *unmarshal) setField(object reflect.Value, value string) error {
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

	fieldName := byts[start+1 : end+1]

	return fieldName, byts[start+end+startNext+1:], nil
}

func (u *unmarshal) getJsonValue(byts []byte) (value []byte, nextValue []byte, err error) {
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

		if numJsonStart == 0 || exit {
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

	return byts[start:end], byts[nextEnd:], nil
}

func (u *unmarshal) getJsonSliceValues(byts []byte) (values [][]byte, err error) {

	var item []byte

	for len(byts) > 0 {
		item, byts, err = u.getJsonValue(byts)
		if err != nil {
			return nil, err
		}
		values = append(values, item)

		if len(byts) > 0 {
			if string(byts[0]) == comma {
				byts = byts[1:]
			}
		}
	}

	return values, nil
}

func (u *unmarshal) getJsonMapValues(byts []byte) (_ map[string][]byte, err error) {

	var key []byte
	var value []byte
	var values = make(map[string][]byte)

	if string(byts[0]) == jsonStart && string(byts[len(byts)-1]) == jsonEnd {
		byts = byts[1 : len(byts)-1]
	}

	for len(byts) > 0 {
		key, byts, err = u.getJsonValue(byts)
		if err != nil {
			return nil, err
		}

		value, byts, err = u.getJsonValue(byts)
		if err != nil {
			return nil, err
		}

		values[string(key)] = value

		exit := false
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

	return values, nil
}

func (u *unmarshal) loadTag(value reflect.Value, typ reflect.StructField) (exists bool, tag string, err error) {
	for _, searchTag := range u.tags {
		tag, exists = typ.Tag.Lookup(searchTag)
		if exists {
			break
		}
	}

	return exists, tag, err
}

func (u *unmarshal) decodeString(str string) string {
	return strings.Replace(str, stringStartEndScaped, stringStartEnd, -1)
}
