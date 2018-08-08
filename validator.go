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

	if v.handlersPre != nil {
		for key, _ := range v.handlersPre {
			handlers[key] = true
		}
	}

	if v.handlersMiddle != nil {
		for key, _ := range v.handlersMiddle {
			handlers[key] = true
		}
	}

	if v.handlersPos != nil {
		for key, _ := range v.handlersPos {
			handlers[key] = true
		}
	}

	return handlers
}

func (v *Validator) AddPre(name string, handler PreTagHandler) *Validator {
	v.handlersPre[name] = handler
	v.activeHandlers[name] = true

	return v
}

func (v *Validator) AddMiddle(name string, handler MiddleTagHandler) *Validator {
	v.handlersMiddle[name] = handler
	v.activeHandlers[name] = true

	return v
}

func (v *Validator) AddPos(name string, handler PosTagHandler) *Validator {
	v.handlersPos[name] = handler
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
