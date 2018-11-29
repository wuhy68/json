package validator

func (v *Validator) NewDefaultMiddleHandlers() map[string]MiddleTagHandler {
	return map[string]MiddleTagHandler{
		ConstTagValue:    v.validate_value,
		ConstTagNot:      v.validate_not,
		ConstTagOptions:  v.validate_options,
		ConstTagSize:     v.validate_size,
		ConstTagMin:      v.validate_min,
		ConstTagMax:      v.validate_max,
		ConstTagNonzero:  v.validate_nonzero,
		ConstTagIszero:  v.validate_iszero,
		ConstTagRegex:    v.validate_regex,
		ConstTagSpecial:  v.validate_special,
		ConstTagSanitize: v.validate_sanitize,
		ConstTagCallback: v.validate_callback,
		ConstTagMatch:    v.validate_match,
		ConstTagSet:      v.validate_set,
		ConstTagDistinct: v.validate_distinct,
	}
}
