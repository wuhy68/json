package validator

import (
	"reflect"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
)

func NewValidator() *Validator {
	return &Validator{
		tag:      "validate",
		handlers: NewDefaultHandlers(),
		log:      logger.NewLogDefault("validator", logger.InfoLevel),
	}
}

type Validator struct {
	tag      string
	handlers map[string]TagHandler
	log      logger.ILogger
}

type TagHandler func(name string, value reflect.Value, expected interface{}) error

// Add ...
func Add(name string, handler TagHandler) (err error) {
	if _, ok := validator.handlers[name]; !ok {
		validator.handlers[name] = handler
	} else {
		err = errors.New("the tag already exists!")
	}

	return err
}

// Validate ...
func Validate(obj interface{}) (err error) {
	return handleValidation(obj)
}
