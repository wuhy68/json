# validator
[![Build Status](https://travis-ci.org/joaosoft/validator.svg?branch=master)](https://travis-ci.org/joaosoft/validator) | [![codecov](https://codecov.io/gh/joaosoft/validator/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/validator) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/validator)](https://goreportcard.com/report/github.com/joaosoft/validator) | [![GoDoc](https://godoc.org/github.com/joaosoft/validator?status.svg)](https://godoc.org/github.com/joaosoft/validator)

A simple struct validator by tags (exported fields only).

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* value (equal to)
* not (not equal to)
* options (one of thous)
* size (size equal to)
* min 
* max 
* nonzero (also supports uuid zero validation)
* regex
* special ( {YYYYMMDD}, {DDMMYYYY}, {date}, {time} )
* error

## With methods for
* AddBefore (add a before-validation)
* AddMiddle (add a middle-validation [by default has all validations])
* AddAfter (add a after-validation [by default has error validation])
* SetErrorCodeHandler (function to get the error when defined with error={xpto:arg1;arg2})
* SetValidateAll (when activated, validates all object instead of stopping on the first error)
* SetTag (set validation tag to other that you define)
* Validate (validate the object)

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/validator
```

## Usage 
This examples are available in the project at [validator/examples](https://github.com/joaosoft/validator/tree/master/examples)

### Code
```go
type Data string

type Example struct {
	Name       string         `validate:"value=joao, dummy, error={1:a;b}, max=10"`
	Age        int            `validate:"value=30, error={99}"`
	Street     int            `validate:"max=10, error=3"`
	Brothers   []Example      `validate:"size=1, error=4"`
	Id         uuid.UUID      `validate:"nonzero, error=5"`
	Option1    string         `validate:"options=aa;bb;cc, error=6"`
	Option2    int            `validate:"options=11;22;33, error=7"`
	Option3    []string       `validate:"options=aa;bb;cc, error=8"`
	Option4    []int          `validate:"options=11;22;33, error=9"`
	Map1       map[string]int `validate:"options=aa:11;bb:22;cc:33, error=10"`
	Map2       map[int]string `validate:"options=11:aa;22:bb;33:cc, error=11"`
	StartTime  string         `validate:"special={time}, error=12"`
	StartDate1 string         `validate:"special={date}, error=13"`
	StartDate2 string         `validate:"special={YYYYMMDD}, error=14"`
	DateString *string        `validate:"special={YYYYMMDD}, error=15"`
	Data       *Data          `validate:"special={YYYYMMDD}, error=16"`
	unexported string
}

var dummy_middle_handler = func(name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	var rtnErrs []error

	err := errors.New("dummy responding...")
	rtnErrs = append(rtnErrs, err)

	return rtnErrs
}

func init() {
	validator.
		AddMiddle("dummy", dummy_middle_handler).
		SetValidateAll(true).
		SetErrorCodeHandler(dummy_error_handler)
}

var errs = map[string]error{
	"1":  errors.New("error 1: a:%s, b:%s"),
	"2":  errors.New("error 2"),
	"3":  errors.New("error 3"),
	"4":  errors.New("error 4"),
	"5":  errors.New("error 5"),
	"6":  errors.New("error 6"),
	"7":  errors.New("error 7"),
	"8":  errors.New("error 8"),
	"9":  errors.New("error 9"),
	"10": errors.New("error 10"),
	"11": errors.New("error 11"),
	"12": errors.New("error 12"),
	"13": errors.New("error 13"),
	"14": errors.New("error 14"),
	"15": errors.New("error 15"),
	"16": errors.New("error 16"),
}
var dummy_error_handler = func(code string, arguments []interface{}, name string, value reflect.Value, expected interface{}, err *[]error) error {
	if err, ok := errs[code]; ok {
		err = fmt.Errorf(err.Error(), arguments...)

		if strings.Contains(err.Error(), "%s") {
			err = fmt.Errorf(err.Error(), name)
		}
		return err
	}
	return nil
}

func main() {
	id, _ := uuid.NewV4()
	str := "2018-12-1"
	data := Data("2018-12-1")
	example := Example{
		Id:         id,
		Name:       "joao",
		Age:        30,
		Street:     10,
		Option1:    "aa",
		Option2:    11,
		Option3:    []string{"aa", "bb", "cc"},
		Option4:    []int{11, 22, 33},
		Map1:       map[string]int{"aa": 11, "bb": 22, "cc": 33},
		Map2:       map[int]string{11: "aa", 22: "bb", 33: "cc"},
		StartTime:  "12:01:00",
		StartDate1: "01-12-2018",
		StartDate2: "2018-12-1",
		DateString: &str,
		Data:       &data,
		Brothers: []Example{
			Example{
				Name:       "jessica",
				Age:        10,
				Street:     12,
				Option1:    "xx",
				Option2:    99,
				Option3:    []string{"aa", "zz", "cc"},
				Option4:    []int{11, 44, 33},
				Map1:       map[string]int{"aa": 11, "kk": 22, "cc": 33},
				Map2:       map[int]string{11: "aa", 22: "bb", 99: "cc"},
				StartTime:  "99:01:00",
				StartDate1: "01-99-2018",
				StartDate2: "2018-99-1",
			},
		},
	}
	if errs := validator.Validate(example); len(errs) > 0 {
		fmt.Printf("ERRORS: %d\n", len(errs))
		for _, err := range errs {
			fmt.Printf("\nERROR: %s", err)
		}
	}
}
```

> ##### Response:
```go
ERRORS: 15

ERROR: error 1: a:a, b:b
ERROR: error 1: a:a, b:b
ERROR: the value [10] is different of the expected [30] on field [Age]
ERROR: {"code":"3","message":"the length [12] is bigger then the expected [10] on field [Street]"}
ERROR: {"code":"4","message":"the length [0] is lower then the expected [1] on field [Brothers]"}
ERROR: {"code":"5","message":"the value shouldn't be zero on field [Id]"}
ERROR: {"code":"6","message":"the value [xx] is different of the expected options [aa;bb;cc] on field [Option1]"}
ERROR: {"code":"7","message":"the value [99] is different of the expected options [11;22;33] on field [Option2]"}
ERROR: {"code":"8","message":"the value [zz] is different of the expected options [aa;bb;cc] on field [Option3]"}
ERROR: {"code":"9","message":"the value [44] is different of the expected options [11;22;33] on field [Option4]"}
ERROR: {"code":"10","message":"the value [22] is different of the expected options [aa:11;bb:22;cc:33] on field [Map1]"}
ERROR: {"code":"11","message":"the value [cc] is different of the expected options [11:aa;22:bb;33:cc] on field [Map2]"}
ERROR: {"code":"12","message":"invalid data [99:01:00] on field [StartTime] "}
ERROR: {"code":"13","message":"invalid data [01-99-2018] on field [StartDate1] "}
ERROR: {"code":"14","message":"invalid data [2018-99-1] on field [StartDate2] "}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
