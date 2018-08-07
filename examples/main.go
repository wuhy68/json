package main

import (
	"fmt"
	"reflect"
	"validator"

	"github.com/joaosoft/errors"
	"github.com/satori/go.uuid"
)

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
