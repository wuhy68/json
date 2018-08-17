package validator

import (
	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
)

func NewValidator() *Validator {

	v := &Validator{
		tag: "validate",
		log: logger.NewLogDefault("validator", logger.InfoLevel),
	}

	v.init()

	return v
}

func (v *Validator) NewActiveHandlers() map[string]bool {
	handlers := make(map[string]bool)

	if v.handlersBefore != nil {
		for key, _ := range v.handlersBefore {
			handlers[key] = true
		}
	}

	if v.handlersMiddle != nil {
		for key, _ := range v.handlersMiddle {
			handlers[key] = true
		}
	}

	if v.handlersAfter != nil {
		for key, _ := range v.handlersAfter {
			handlers[key] = true
		}
	}

	return handlers
}

func (v *Validator) AddBefore(name string, handler BeforeTagHandler) *Validator {
	v.handlersBefore[name] = handler
	v.activeHandlers[name] = true

	return v
}

func (v *Validator) AddMiddle(name string, handler MiddleTagHandler) *Validator {
	v.handlersMiddle[name] = handler
	v.activeHandlers[name] = true

	return v
}

func (v *Validator) AddAfter(name string, handler AfterTagHandler) *Validator {
	v.handlersAfter[name] = handler
	v.activeHandlers[name] = true

	return v
}

func (v *Validator) SetErrorCodeHandler(handler ErrorCodeHandler) *Validator {
	v.errorCodeHandler = handler

	return v
}

func (v *Validator) SetValidateAll(validateAll bool) *Validator {
	v.validateAll = validateAll

	return v
}

func (v *Validator) SetTag(tag string) *Validator {
	v.tag = tag

	return v
}

// Validate ...
func (v *Validator) Validate(obj interface{}) *errors.ListErr {
	return handleValidation(obj)
}
