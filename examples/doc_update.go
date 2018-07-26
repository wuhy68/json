package main

import (
	"fmt"

	"strconv"

	log "github.com/joaosoft/logger"
)

func updateDocumentWithId(id string) {

	// document update with id
	age, _ := strconv.Atoi(id)
	id, err := client.Create().Index("persons").Type("person").Id(id).Body(person{
		Name: "luis",
		Age:  age + 20,
	}).Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\nupdated person with id %s\n", id)
	}
}
