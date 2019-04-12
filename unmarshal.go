package json

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
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
	index := 0
	switch string(byts[index : index+1]) {
	case jsonStart:
		index++

		start := bytes.Index(byts[index:], []byte(stringStartEnd))
		end := bytes.Index(byts[index+start+1:], []byte(stringStartEnd))

		fieldName := string(byts[index+start+1 : index+end+1])
		fmt.Println("Field Name:" + fieldName)

		exists, field, err := u.getField(object, fieldName)
		if err != nil {
			return err
		}

		if exists {
			u.setField(field, "nabo")
		}

	case arrayStart:

	default:
	}
	return nil
}

func (u *unmarshal) getField(object reflect.Value, name string) (bool, reflect.Value, error) {
	hasImplementation := object.Type().Implements(typeUnmarshal)
	isSlice := object.Kind() == reflect.Slice && !hasImplementation
	isMap := object.Kind() == reflect.Map && !hasImplementation
	isMapOfSlices := isMap && object.Type().Elem().Kind() == reflect.Slice

	if isMapOfSlices {
		object = reflectAlloc(object.Type().Elem().Elem())
	} else if isSlice || isMap {
		object = reflectAlloc(object.Type().Elem())
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
	case reflect.String:
		object.SetString(value)
	case reflect.Int:
		tmpValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		object.SetInt(tmpValue)
	}

	return nil
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
