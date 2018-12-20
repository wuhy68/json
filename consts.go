package validator

const (
	ConstTagForDateDefault  = "date"
	ConstTagForDateYYYYMMDD = "YYYYMMDD"
	ConstTagForDateDDMMYYYY = "DDMMYYYY"
	ConstTagForTimeDefault  = "time"
	ConstTagForTimeHHMMSS   = "HHMMSS"
	ConstTagForURL          = "url"
	ConstTagForEmail        = "email"

	ConstRegexForDateDDMMYYYY = `^(0?[1-9]|[12][0-9]|3[01])(/|-|.)([1-9]|0[0-9]|1[0-2])(/|-|.)[0-9]{4}$`
	ConstRegexForDateYYYYMMDD = `^[0-9]{4}(/|-|.)(0?[1-9]|[12][0-9]|3[01])(/|-|.)([1-9]|0[0-9]|1[0-2])$`
	ConstRegexForDateDefault  = ConstRegexForDateDDMMYYYY
	ConstRegexForTimeDefault  = ConstRegexForTimeHHMMSS
	ConstRegexForTimeHHMMSS   = `^([0-9]|0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`
	ConstRegexForErrorTag     = "{[A-Za-z0-9_-]+:?([A-Za-z0-9_-];?)+}"
	ConstRegexForURL          = "^((http|https)://)?(www)?[a-zA-Z0-9-._:/?&=,]+$"
	ConstRegexForEmail        = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

	ConstDefaultValidationTag = "validate"
	ConstDefaultLogTag        = "validator"

	ConstTagId       = "id"
	ConstTagValue    = "value"
	ConstTagError    = "error"
	ConstTagIf       = "if"
	ConstTagNot      = "not"
	ConstTagOptions  = "options"
	ConstTagSize     = "size"
	ConstTagMin      = "min"
	ConstTagMax      = "max"
	ConstTagNonzero  = "nonzero"
	ConstTagIszero   = "iszero"
	ConstTagRegex    = "regex"
	ConstTagSpecial  = "special"
	ConstTagSanitize = "sanitize"
	ConstTagCallback = "callback"
	ConstTagSet      = "set"
	ConstTagTrim     = "trim"
	ConstTagKey      = "key"
	ConstTagDistinct = "distinct"
	ConstTagAlpha    = "alpha"
	ConstTagNumeric  = "numeric"
)
