package main

import (
	"fmt"
	"go-error/service"
)

func main() {

	err := goerror.NewError(fmt.Errorf("erro 1"))
	err.Add(fmt.Errorf("erro 2"))
	err.Add(fmt.Errorf("erro 3"))

	fmt.Printf("Error: %s, Cause: %s", err.Error(), err.Cause())
}
