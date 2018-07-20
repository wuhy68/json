package main

import (
	"elastic"

	log "github.com/joaosoft/logger"
)

func main() {
	var data []interface{}

	client := elastic.NewElastic()
	// you can define the configuration without having a configuration file
	//client1 := elastic.NewElastic(elastic.WithConfiguration(elastic.NewConfig("http://localhost:9200")))

	d1 := elastic.TemplateData{Data: map[string]interface{}{"default_plan": true}}

	err := client.Search().
		Index("persons").
		Document("person").
		Object(&data).
		Template("/examples/templates", "get.example.1.template", &d1, false).
		Execute()

	if err != nil {
		log.Error(err)
	}

	d2 := elastic.TemplateData{Data: 123}
	err = client.Search().
		Index("persons").
		Document("person").
		Object(&data).
		Template("/examples/templates", "get.example.2.template", &d2, false).
		Execute()

	if err != nil {
		log.Error(err)
	}
}
