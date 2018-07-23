package main

import (
	"elastic"

	"fmt"

	"os"

	log "github.com/joaosoft/logger"
)

func searchDocument(name string) {
	var data []person

	client := elastic.NewElastic()
	// you can define the configuration without having a configuration file
	//client1 := elastic.NewElastic(elastic.WithConfiguration(elastic.NewConfig("http://localhost:9200")))

	d1 := elastic.TemplateData{Data: map[string]interface{}{"name": name}}

	// document search
	dir, _ := os.Getwd()
	err := client.Search().
		Index("persons").
		Type("person").
		Object(&data).
		Template(dir+"/examples/templates", "get.example.1.template", &d1, false).
		Execute()

	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("\nsearch person by name:%s %+v\n", name, data)
	}
}
