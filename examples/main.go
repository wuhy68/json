package main

import (
	"fmt"
	"reflect"
	"validator"

	"github.com/joaosoft/errors"
	"github.com/satori/go.uuid"
)

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
			fmt.Printf("\nCODE: %s, MESSAGE: %s", err.GetCode(), err.GetError())
		}
	}
}
