# validator
[![Build Status](https://travis-ci.org/joaosoft/validator.svg?branch=master)](https://travis-ci.org/joaosoft/validator) | [![codecov](https://codecov.io/gh/joaosoft/validator/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/validator) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/validator)](https://goreportcard.com/report/github.com/joaosoft/validator) | [![GoDoc](https://godoc.org/github.com/joaosoft/validator?status.svg)](https://godoc.org/github.com/joaosoft/validator)

A simple struct validator by tags.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* value
* options
* size
* min 
* max 
* nonzero (also supports uuid zero validation)
* regex
* error

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
type Example struct {
	Name     string         `validate:"value=joao, dummy, error={1}, max=10"`
	Age      int            `validate:"value=30, error=2"`
	Street   int            `validate:"max=10, error=3"`
	Brothers []Example      `validate:"size=1, error=4"`
	Id       uuid.UUID      `validate:"nonzero, error=5"`
	Option1  string         `validate:"options=aa;bb;cc, error=6"`
	Option2  int            `validate:"options=11;22;33, error=7"`
	Option3  []string       `validate:"options=aa;bb;cc, error=8"`
	Option4  []int          `validate:"options=11;22;33, error=9"`
	Map1     map[string]int `validate:"options=aa:11;bb:22;cc:33, error=10"`
	Map2     map[int]string `validate:"options=11:aa;22:bb;33:cc, error=11"`
}

var dummy_middle_handler = func(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
	return errors.New("0", "dummy responding...")
}

func init() {
	validator.AddMiddle("dummy", dummy_middle_handler).SetValidateAll(true).SetErrorCodeHandler(dummy_error_handler)
}

var errs = map[string]errors.IErr{
	"1":  errors.New("1", "Error 1"),
	"2":  errors.New("2", "Error 2"),
	"3":  errors.New("3", "Error 3"),
	"4":  errors.New("4", "Error 4"),
	"5":  errors.New("5", "Error 5"),
	"6":  errors.New("6", "Error 6"),
	"7":  errors.New("7", "Error 7"),
	"8":  errors.New("8", "Error 8"),
	"9":  errors.New("9", "Error 9"),
	"10": errors.New("10", "Error 10"),
	"11": errors.New("11", "Error 11"),
}
var dummy_error_handler = func(code string) errors.IErr {
	return errs[code]
}

func main() {
	id, _ := uuid.NewV4()
	example := Example{
		Id:      id,
		Name:    "joao",
		Age:     30,
		Street:  10,
		Option1: "aa",
		Option2: 11,
		Option3: []string{"aa", "bb", "cc"},
		Option4: []int{11, 22, 33},
		Map1:    map[string]int{"aa": 11, "bb": 22, "cc": 33},
		Map2:    map[int]string{11: "aa", 22: "bb", 33: "cc"},
		Brothers: []Example{
			Example{
				Name:    "jessica",
				Age:     10,
				Street:  12,
				Option1: "xx",
				Option2: 99,
				Option3: []string{"aa", "zz", "cc"},
				Option4: []int{11, 44, 33},
				Map1:    map[string]int{"aa": 11, "kk": 22, "cc": 33},
				Map2:    map[int]string{11: "aa", 22: "bb", 99: "cc"},
			},
		},
	}
	if e := validator.Validate(example); e.Len() > 0 {
		fmt.Printf("ERRORS: %d\n", e.Len())
		for _, err := range *e {
			fmt.Printf("\nCODE: %s, MESSAGE: %s", err.GetCode(), err.GetError())
		}
	}
}
```

> ##### Response:
```go
ERRORS: 12

CODE: 1, MESSAGE: Error 1
CODE: 1, MESSAGE: Error 1
CODE: 2, MESSAGE: the value [10] is diferent of the expected [30] on field [Age]
CODE: 3, MESSAGE: the length [12] is bigger then the expected [10] on field [Street]
CODE: 4, MESSAGE: the length [0] is lower then the expected [1] on field [Brothers]
CODE: 5, MESSAGE: the field [Id] shouldn't be zero
CODE: 6, MESSAGE: the value [xx] is diferent of the expected options [aa;bb;cc] on field [Option1]
CODE: 7, MESSAGE: the value [99] is diferent of the expected options [11;22;33] on field [Option2]
CODE: 8, MESSAGE: the value [zz] is diferent of the expected options [aa;bb;cc] on field [Option3]
CODE: 9, MESSAGE: the value [44] is diferent of the expected options [11;22;33] on field [Option4]
CODE: 10, MESSAGE: the value [kk: 22] is diferent of the expected options [aa:11;bb:22;cc:33] on field [Map1]
CODE: 11, MESSAGE: the value [99: cc] is diferent of the expected options [11:aa;22:bb;33:cc] on field [Map2]
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
