package main

import (
	"fmt"

	log "github.com/joaosoft/logger"
)

func deleteIndex() {

	err := client.DeleteIndex().Index("persons").Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\ndeleted persons index\n")
	}
}
