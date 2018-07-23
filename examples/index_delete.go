package main

import (
	"elastic"
	"fmt"

	log "github.com/joaosoft/logger"
)

func deleteIndex() {
	client := elastic.NewElastic()
	// you can define the configuration without having a configuration file
	//client1 := elastic.NewElastic(elastic.WithConfiguration(elastic.NewConfig("http://localhost:9200")))

	err := client.DeleteIndex().Index("persons").Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\ndeleted persons index\n")
	}
}
