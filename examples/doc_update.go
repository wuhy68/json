package main

import (
	"elastic"

	"fmt"

	"strconv"

	log "github.com/joaosoft/logger"
)

func updateDocumentWithId(id string) {
	client := elastic.NewElastic()
	// you can define the configuration without having a configuration file
	//client1 := elastic.NewElastic(elastic.WithConfiguration(elastic.NewConfig("http://localhost:9200")))

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
