package main

import (
	"fmt"

	"strconv"

	log "github.com/joaosoft/logger"
)

func createDocumentWithId(id string) {

	// document create with id
	age, _ := strconv.Atoi(id)
	id, err := client.Create().Index("persons").Type("person").Id(id).Body(person{
		Name: "joao",
		Age:  age + 20,
	}).Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\ncreated a new person with id %s\n", id)
	}
}

func createDocumentWithoutId() string {

	// document create without id
	id, err := client.Create().Index("persons").Type("person").Body(person{
		Name: "joao",
		Age:  30,
	}).Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\ncreated a new person with id %s\n", id)
	}

	return id
}
