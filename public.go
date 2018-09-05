package validator

func AddBefore(name string, handler BeforeTagHandler) *Validator {
	return validatorInstance.AddBefore(name, handler)
}

func AddMiddle(name string, handler MiddleTagHandler) *Validator {
	return validatorInstance.AddMiddle(name, handler)
}

func AddAfter(name string, handler AfterTagHandler) *Validator {
	return validatorInstance.AddAfter(name, handler)
}

func SetValidateAll(validate bool) *Validator {
	return validatorInstance.SetValidateAll(validate)
}

func SetTag(tag string) *Validator {
	return validatorInstance.SetTag(tag)
}

// Validate ...
func Validate(obj interface{}) []error {
	return validatorInstance.Validate(obj)
}
