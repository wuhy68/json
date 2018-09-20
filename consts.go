package validator

const (
	TagForDateDefault  = "date"
	TagForDateYYYYMMDD = "YYYYMMDD"
	TagForDateDDMMYYYY = "DDMMYYYY"
	TagForTimeDefault  = "time"
	TagForTimeHHMMSS   = "HHMMSS"
	TagForURL          = "url"

	RegexForDateDDMMYYYY = `^(0?[1-9]|[12][0-9]|3[01])(/|-|.)([1-9]|0[0-9]|1[0-2])(/|-|.)[0-9]{4}$`
	RegexForDateYYYYMMDD = `^[0-9]{4}(/|-|.)(0?[1-9]|[12][0-9]|3[01])(/|-|.)([1-9]|0[0-9]|1[0-2])$`
	RegexForDateDefault  = RegexForDateDDMMYYYY
	RegexForTimeDefault  = RegexForTimeHHMMSS
	RegexForTimeHHMMSS   = `^([0-9]|0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`
	RegexForErrorTag     = "{[A-Za-z0-9_-]+:?([A-Za-z0-9_-];?)+}"
	RegexForURL          = "^((http|https)://)?(www)?[a-zA-Z0-9-._:/?&=,]+$"
)
