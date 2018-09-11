package validator

func (v *Validator) NewDefaultPosHandlers() map[string]AfterTagHandler {
	return map[string]AfterTagHandler{
		"error": v.validate_error,
	}
}
