package validator

import (
	"reflect"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
)

func (v *Validator) init() {
	v.handlersPre = v.NewDefaultPreHandlers()
	v.handlersMiddle = v.NewDefaultMiddleHandlers()
	v.handlersPos = v.NewDefaultPosHandlers()
	v.activeHandlers = v.NewActiveHandlers()

}

type Validator struct {
	tag              string
	activeHandlers   map[string]bool
	handlersPre      map[string]PreTagHandler
	handlersMiddle   map[string]MiddleTagHandler
	handlersPos      map[string]PosTagHandler
	errorCodeHandler ErrorCodeHandler
	log              logger.ILogger
	validateAll      bool
}

type ErrorCodeHandler func(code string) error

type PreTagHandler func(name string, value reflect.Value, expected interface{}) errors.ListErr
type MiddleTagHandler func(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.ListErr
type PosTagHandler func(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.ListErr
