package main

import (
	"elastic"

	"fmt"

	"os"

	log "github.com/joaosoft/logger"
)

func countOnIndex(name string) int64 {

	d1 := elastic.CountTemplate{Data: map[string]interface{}{"name": name}}

	// index count
	dir, _ := os.Getwd()
	count, err := client.Count().
		Index("persons").
		Template(dir+"/examples/templates", "get.example.count.template", &d1, false).
		Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\ncount persons with name %s: %d\n", name, count)
	}

	return count
}

func countOnDocument(name string) int64 {

	d1 := elastic.CountTemplate{Data: map[string]interface{}{"name": name}}

	// index count
	dir, _ := os.Getwd()
	count, err := client.Count().
		Index("persons").
		Type("person").
		Template(dir+"/examples/templates", "get.example.count.template", &d1, false).
		Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\ncount persons with name %s: %d\n", name, count)
	}

	return count
}
