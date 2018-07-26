package main

import (
	"elastic"

	"fmt"

	"os"

	log "github.com/joaosoft/logger"
)

func searchDocument(name string) {
	var data []person

	d1 := elastic.SearchTemplate{Data: map[string]interface{}{"name": name}}

	// document search
	dir, _ := os.Getwd()
	err := client.Search().
		Index("persons").
		Type("person").
		Object(&data).
		Template(dir+"/examples/templates", "get.example.search.template", &d1, false).
		Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\nsearch person by name:%s %+v\n", name, data)
	}
}
