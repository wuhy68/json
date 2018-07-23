package main

import (
	"elastic"

	"fmt"

	"strconv"

	log "github.com/joaosoft/logger"
)

func createScroll(id string) {
	client := elastic.NewElastic()
	// you can define the configuration without having a configuration file
	//client1 := elastic.NewElastic(elastic.WithConfiguration(elastic.NewConfig("http://localhost:9200")))

	// document create with id
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
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
