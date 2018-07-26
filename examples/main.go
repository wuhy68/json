package main

import (
	"elastic"
	"time"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var client = elastic.NewElastic()

// you can define the configuration without having a configuration file
//client1 := elastic.NewElastic(elastic.WithConfiguration(elastic.NewConfig("http://localhost:9200")))

func main() {

	// index create with mapping
	createIndexWithMapping()

	// document create
	createDocumentWithId("1")
	createDocumentWithId("2")
	generatedId := createDocumentWithoutId()

	// document update
	updateDocumentWithId("1")
	updateDocumentWithId("2")

	// document search
	// wait elastic to index the last update...
	<-time.After(time.Second * 2)
	searchDocument("luis")

	// count index documents
	countOnIndex("luis")
	countOnDocument("luis")

	// document delete
	deleteDocumentWithId(generatedId)

	// index exists
	existsIndex("persons")
	existsIndex("bananas	")

	// index delete
	deleteIndex()
}
