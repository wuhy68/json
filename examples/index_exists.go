package main

import (
	"elastic"
	"fmt"

	"net/http"

	log "github.com/joaosoft/logger"
)

func existsIndex(index string) {
	client := elastic.NewElastic()
	// you can define the configuration without having a configuration file
	//client1 := elastic.NewElastic(elastic.WithConfiguration(elastic.NewConfig("http://localhost:9200")))

	status, err := client.ExistsIndex().Index(index).Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\nexists index? %t\n", status == http.StatusOK)
	}
}
