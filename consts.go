package json

import "reflect"

const (
	jsonTag              = "json"
	jsonStart            = "{"
	jsonEnd              = "}"
	arrayStart           = "["
	arrayEnd             = "]"
	stringStartEnd       = `"`
	stringStartEndScaped = `\"`
	is                   = `:`
	comma                = ","
	null                 = "null"
	empty                = " "
	booleanTrue          = "true"
	booleanFalse         = "false"
)

var (
	TypeString = reflect.TypeOf("")
	TypeFloat64 = reflect.TypeOf(0.0)
	TypeBoolean = reflect.TypeOf(true)
)