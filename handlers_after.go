package validator

func (v *Validator) NewDefaultPosHandlers() map[string]AfterTagHandler {
	return map[string]AfterTagHandler{
		ConstTagError: v.validate_error,
	}
}
