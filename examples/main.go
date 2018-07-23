package main

import "time"

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

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

	// document delete
	deleteDocumentWithId(generatedId)

	// index exists
	existsIndex("persons")
	existsIndex("bananas	")

	// index delete
	deleteIndex()
}
