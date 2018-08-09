package main

import (
	"fmt"
	"reflect"
	"validator"

	"github.com/joaosoft/errors"
	"github.com/satori/go.uuid"
)

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

var dummy_middle_handler = func(name string, value reflect.Value, expected interface{}, errs *errors.ListErr) errors.ListErr {
	rtnErrs := make(errors.ListErr, 0)

	rtnErrs = append(rtnErrs, errors.New("0", "dummy responding..."))

	return rtnErrs
}

func init() {
	validator.
		AddMiddle("dummy", dummy_middle_handler).
		SetValidateAll(true).
		SetErrorCodeHandler(dummy_error_handler)
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
	if e := validator.Validate(example); !e.IsEmpty() {
		fmt.Printf("ERRORS: %d\n", e.Len())
		for _, err := range *e {
			fmt.Printf("\nCODE: %s, MESSAGE: %s", err.GetCode(), err.GetError())
		}
	}
}
