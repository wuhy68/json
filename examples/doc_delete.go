package main

import (
	"elastic"

	"fmt"

	log "github.com/joaosoft/logger"
)

func deleteDocumentWithId(id string) {
	client := elastic.NewElastic()
	// you can define the configuration without having a configuration file
	//client1 := elastic.NewElastic(elastic.WithConfiguration(elastic.NewConfig("http://localhost:9200")))

	err := client.Delete().Index("persons").Type("person").Id(id).Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\ndeleted person with id %s\n", "1")
	}
}