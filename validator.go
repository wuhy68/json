package validator

import (
	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
)

func NewValidator() *Validator {

	preHandlers := NewDefaultPreHandlers()
	middleHandlers := NewDefaultMiddleHandlers()
	posHandlers := NewDefaultPosHandlers()

	return &Validator{
		tag:            "validate",
		handlersPre:    preHandlers,
		handlersMiddle: middleHandlers,
		handlersPos:    posHandlers,
		activeHandlers: loadActiveHandlers(preHandlers, middleHandlers, posHandlers),
		log:            logger.NewLogDefault("validator", logger.InfoLevel),
	}
}

func loadActiveHandlers(preHandlers map[string]PreTagHandler, middleHandlers map[string]MiddleTagHandler, posHandlers map[string]PosTagHandler) map[string]bool {
	handlers := make(map[string]bool)

	for key, _ := range preHandlers {
		handlers[key] = true
	}

	for key, _ := range middleHandlers {
		handlers[key] = true
	}

	for key, _ := range posHandlers {
		handlers[key] = true
	}

	return handlers
}

func (v *Validator) AddPre(name string, handler PreTagHandler) *Validator {
	v.handlersPre[name] = handler
	v.activeHandlers[name] = true

	return validator
}

func (v *Validator) AddMiddle(name string, handler MiddleTagHandler) *Validator {
	v.handlersMiddle[name] = handler
	v.activeHandlers[name] = true

	return validator
}

func (v *Validator) AddPos(name string, handler PosTagHandler) *Validator {
	v.handlersPos[name] = handler
	v.activeHandlers[name] = true

	return validator
}

func (v *Validator) SetValidateAll(validateAll bool) {
	v.validateAll = validateAll
}

func (v *Validator) SetTag(tag string) {
	v.tag = tag
}

// Validate ...
func (v *Validator) Validate(obj interface{}) *errors.ListErr {
	return handleValidation(obj)
}
