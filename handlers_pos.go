package validator

func (v *Validator) NewDefaultPosHandlers() map[string]PosTagHandler {
	return map[string]PosTagHandler{"error": v.validate_error}
}