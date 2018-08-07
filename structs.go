package validator

import (
	"reflect"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
)

type Validator struct {
	tag            string
	activeHandlers map[string]bool
	handlersPre    map[string]PreTagHandler
	handlersMiddle map[string]MiddleTagHandler
	handlersPos    map[string]PosTagHandler
	log            logger.ILogger
	validateAll    bool
}

type PreTagHandler func(name string, value reflect.Value, expected interface{}) *errors.Err
type MiddleTagHandler func(name string, value reflect.Value, expected interface{}, err *errors.ListErr) *errors.Err
type PosTagHandler func(name string, value reflect.Value, expected interface{}, err *errors.ListErr) *errors.Err
