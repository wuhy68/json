package main

import (
	"fmt"

	log "github.com/joaosoft/logger"
)

func deleteDocumentWithId(id string) {

	err := client.Delete().Index("persons").Type("person").Id(id).Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\ndeleted person with id %s\n", "1")
	}
}