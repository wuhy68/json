# validator
[![Build Status](https://travis-ci.org/joaosoft/validator.svg?branch=master)](https://travis-ci.org/joaosoft/validator) | [![codecov](https://codecov.io/gh/joaosoft/validator/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/validator) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/validator)](https://goreportcard.com/report/github.com/joaosoft/validator) | [![GoDoc](https://godoc.org/github.com/joaosoft/validator?status.svg)](https://godoc.org/github.com/joaosoft/validator)

A simple struct validator by tags.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* value
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
	Name     string    `validate:"value=joao, dummy, error={1}, max=10"`
	Age      int       `validate:"value=30, error=2"`
	Street   int       `validate:"max=10, error=3"`
	Brothers []Example `validate:"size=1, error=4"`
	Id       uuid.UUID `validate:"nonzero, error=5"`
}

var dummy_middle_handler = func(name string, value reflect.Value, expected interface{}, err *errors.ListErr) errors.IErr {
	return errors.New("0", "dummy responding...")
}

func init() {
	validator.AddMiddle("dummy", dummy_middle_handler).SetValidateAll(true).SetErrorCodeHandler(dummy_error_handler)
}

var errs = map[string]errors.IErr{
	"1": errors.New("1", "Error 1"),
	"2": errors.New("1", "Error 2"),
	"3": errors.New("1", "Error 3"),
	"4": errors.New("1", "Error 4"),
	"5": errors.New("1", "Error 5"),
}
var dummy_error_handler = func(code string) errors.IErr {
	return errs[code]
}

func main() {
	id, _ := uuid.NewV4()
	example := Example{
		Id:     id,
		Name:   "joao",
		Age:    30,
		Street: 10,
		Brothers: []Example{
			Example{
				Name:   "jessica",
				Age:    10,
				Street: 12,
			},
		},
	}
	if e := validator.Validate(example); e.Len() > 0 {
		fmt.Printf("ERRORS: %d\n", e.Len())
		for _, err := range *e {
			fmt.Printf("\nCODE: %s, MESAGE: %s", err.GetCode(), err.GetError())
		}
	}
}
```

> ##### Response:
```go
ERRORS: 7

CODE: 1, MESAGE: Error 1
CODE: 1, MESAGE: Error 1
CODE: 1, MESAGE: Error 1
CODE: 2, MESAGE: the value [10] is diferent of the expected [30] on field [Age]
CODE: 3, MESAGE: the length [12] is bigger then the expected [10] on field [Street]
CODE: 4, MESAGE: the length [0] is lower then the expected [1] on field [Brothers]
CODE: 5, MESAGE: the field [Id] shouldn't be zero
Process finished with exit code 0

```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
