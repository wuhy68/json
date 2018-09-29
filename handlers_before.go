package validator

func (v *Validator) NewDefaultBeforeHandlers() map[string]BeforeTagHandler {
	return map[string]BeforeTagHandler{
		"id": v.validate_id,
		"if": v.validate_if,
	}
}
