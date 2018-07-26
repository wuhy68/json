package main

import (
	"fmt"

	"net/http"

	log "github.com/joaosoft/logger"
)

func existsIndex(index string) {

	status, err := client.ExistsIndex().Index(index).Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\nexists index? %t\n", status == http.StatusOK)
	}
}
