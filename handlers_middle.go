package validator

func (v *Validator) NewDefaultMiddleHandlers() map[string]MiddleTagHandler {
	return map[string]MiddleTagHandler{
		"value":    v.validate_value,
		"not":      v.validate_not,
		"options":  v.validate_options,
		"size":     v.validate_size,
		"min":      v.validate_min,
		"max":      v.validate_max,
		"nonzero":  v.validate_nonzero,
		"regex":    v.validate_regex,
		"special":  v.validate_special,
		"sanitize": v.validate_sanitize,
		"callback": v.validate_callback,
		"match":    v.validate_match,
	}
}
