package main

import (
	"fmt"
	"reflect"
	"validator"
)

type Example struct {
	Name     string    `validate:"value=joao, tagOne=teste"`
	Age      int       `validate:"value=30"`
	Street   int       `validate:"max=10"`
	Brothers []Example `validate:"size=1"`
}

var tagOne_handler = func(name string, value reflect.Value, expected interface{}) error {
	fmt.Printf("hello tagOne!")
	return nil
}

func init() {
	if err := validator.Add("tagOne", tagOne_handler); err != nil {
		fmt.Printf("error adding tag %s", "tagOne")
	}
}

func main() {
	example := Example{
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
	if err := validator.Validate(example); err != nil {
		fmt.Printf("\nvalidation failed with error [%s]", err.Error())
	}
}
