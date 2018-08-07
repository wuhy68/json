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
import "github.com/joaosoft/validator"

type Example struct {
	Name     string    `validate:"value=joao, dummy, error=1, max=10"`
	Age      int       `validate:"value=30, error=2"`
	Street   int       `validate:"max=10, error=3"`
	Brothers []Example `validate:"size=1, error=4"`
	Id       uuid.UUID `validate:"nonzero, error=5"`
}

var dummy_middle_handler = func(name string, value reflect.Value, expected interface{}, err *errors.ListErr) *errors.Err {
	return errors.New("0", "dummy responding...")
}

func init() {
	validator.AddMiddle("dummy", dummy_middle_handler).SetValidateAll(true)
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
	if e := validator.Validate(example); len(*e) > 0 {
		fmt.Printf("ERRORS: %d\n", len(*e))
		for _, err := range *e {
			fmt.Printf("\nCODE: %s, MESAGE: %s", err.Code, err.Err)
		}
	}
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
