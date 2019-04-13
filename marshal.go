package json

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

var typeMarshal = reflect.TypeOf((*imarshal)(nil)).Elem()

type imarshal interface {
	MarshalJSON() ([]byte, error)
}

type marshal struct {
	object interface{}
	result *bytes.Buffer
	tags   []string
}

func newMarshal(object interface{}, tags ...string) *marshal {
	if len(tags) == 0 {
		tags = append(tags, jsonTag)
	}
	return &marshal{object: object, result: bytes.NewBuffer(make([]byte, 0)), tags: tags}
}

func (m *marshal) execute() ([]byte, error) {
	err := m.do(reflect.ValueOf(m.object))
	return m.result.Bytes(), err
}

func (m *marshal) do(object reflect.Value) error {
	hasImplementation := object.Type().Implements(typeMarshal)

	types := reflect.TypeOf(object.Interface())

	if !object.CanInterface() {
		return nil
	}

	if object.Kind() == reflect.Ptr && !object.IsNil() {
		object = object.Elem()

		if object.IsValid() {
			types = object.Type()
		} else {
			return nil
		}
	}

	switch object.Kind() {
	case reflect.Struct:
		if _, err := m.result.WriteString(jsonStart); err != nil {
			return err
		}

		if hasImplementation {
			if val, ok := object.Interface().(imarshal); ok {
				bytes, err := val.MarshalJSON()
				if err != nil {
					return err
				}

				if _, err := m.result.Write(bytes); err != nil {
					return err
				}
				return nil
			}
		}

		addComma := false
		for i := 0; i < types.NumField(); i++ {
			nextValue := object.Field(i)
			nextType := types.Field(i)

			if nextValue.Kind() == reflect.Ptr && !nextValue.IsNil() {
				nextValue = nextValue.Elem()
			}

			if !nextValue.CanInterface() {
				continue
			}

			exists, tag, err := m.loadTag(nextValue, nextType)
			if err != nil {
				return err
			}

			if !exists {
				continue
			}

			if addComma {
				if _, err := m.result.WriteString(comma); err != nil {
					return err
				}
			}

			if _, err := m.result.WriteString(fmt.Sprintf(`%s%s%s%s`, stringStartEnd, tag, stringStartEnd, is)); err != nil {
				return err
			}
			addComma = true

			if err := m.do(nextValue); err != nil {
				return err
			}
		}

		if _, err := m.result.WriteString(jsonEnd); err != nil {
			return err
		}

	case reflect.Array, reflect.Slice:
		if object.IsNil() {
			if _, err := m.result.WriteString(null); err != nil {
				return err
			}
			return nil
		}

		if _, err := m.result.WriteString(arrayStart); err != nil {
			return err
		}

		addComma := false
		for i := 0; i < object.Len(); i++ {
			nextValue := object.Index(i)

			if !nextValue.CanInterface() {
				continue
			}

			if addComma {
				if _, err := m.result.WriteString(comma); err != nil {
					return err
				}
			}

			if err := m.do(nextValue); err != nil {
				return err
			}
			addComma = true
		}

		if _, err := m.result.WriteString(arrayEnd); err != nil {
			return err
		}

	case reflect.Map:
		if object.IsNil() {
			if _, err := m.result.WriteString(null); err != nil {
				return err
			}
			return nil
		}

		if _, err := m.result.WriteString(jsonStart); err != nil {
			return err
		}

		addComma := false
		for _, key := range object.MapKeys() {
			nextValue := object.MapIndex(key)

			if !nextValue.CanInterface() {
				continue
			}

			if addComma {
				if _, err := m.result.WriteString(comma); err != nil {
					return err
				}
			}

			if err := m.handleKey(key); err != nil {
				return err
			}

			if err := m.do(nextValue); err != nil {
				return err
			}
			addComma = true
		}

		if _, err := m.result.WriteString(jsonEnd); err != nil {
			return err
		}

	default:
		if err := m.handleValue(object); err != nil {
			return err
		}
	}
	return nil
}

func (m *marshal) handleKey(key reflect.Value) error {
	switch key.Kind() {
	case reflect.String:
		if _, err := m.result.WriteString(fmt.Sprintf("%s%s", m.encodeString(fmt.Sprintf(`%+v`, key.Interface())), is)); err != nil {
			return err
		}
	default:
		if _, err := m.result.WriteString(fmt.Sprintf(`%+v%s`, key.Interface(), is)); err != nil {
			return err
		}
	}

	return nil
}

func (m *marshal) handleValue(object reflect.Value) error {

	switch object.Kind() {
	case reflect.String:
		if _, err := m.result.WriteString(m.encodeString(fmt.Sprintf(`%+v`, object.Interface()))); err != nil {
			return err
		}
	default:
		if object.Kind() == reflect.Ptr && object.IsNil() {
			if _, err := m.result.WriteString(fmt.Sprintf(`%s`, null)); err != nil {
				return err
			}
			return nil
		}

		if _, err := m.result.WriteString(fmt.Sprintf(`%+v`, object.Interface())); err != nil {
			return err
		}
	}

	return nil
}

func (m *marshal) loadTag(object reflect.Value, typ reflect.StructField) (exists bool, tag string, err error) {
	for _, searchTag := range m.tags {
		tag, exists = typ.Tag.Lookup(searchTag)
		if exists {
			break
		}
	}

	return exists, tag, err
}

func (m *marshal) encodeString(str string) string {
	return fmt.Sprintf("%s%s%s", stringStartEnd, strings.Replace(str, stringStartEnd, stringStartEndScaped, -1), stringStartEnd)
}
