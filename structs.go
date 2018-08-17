package validator

import (
	"reflect"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
)

func (v *Validator) init() {
	v.handlersBefore = v.NewDefaultBeforeHandlers()
	v.handlersMiddle = v.NewDefaultMiddleHandlers()
	v.handlersAfter = v.NewDefaultPosHandlers()
	v.activeHandlers = v.NewActiveHandlers()

}

type Validator struct {
	tag              string
	activeHandlers   map[string]bool
	handlersBefore   map[string]BeforeTagHandler
	handlersMiddle   map[string]MiddleTagHandler
	handlersAfter    map[string]AfterTagHandler
	errorCodeHandler ErrorCodeHandler
	log              logger.ILogger
	validateAll      bool
}

type ErrorCodeHandler func(code string) error

type BeforeTagHandler func(name string, value reflect.Value, expected interface{}) errors.ListErr
type MiddleTagHandler func(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.ListErr
type AfterTagHandler func(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.ListErr
