# elastic
[![Build Status](https://travis-ci.org/joaosoft/elastic.svg?branch=master)](https://travis-ci.org/joaosoft/elastic) | [![codecov](https://codecov.io/gh/joaosoft/elastic/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/elastic) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/elastic)](https://goreportcard.com/report/github.com/joaosoft/elastic) | [![GoDoc](https://godoc.org/github.com/joaosoft/elastic?status.svg)](https://godoc.org/github.com/joaosoft/elastic)

A simple and fast elastic client.

## Support for 
> Create / Exists / Delete index (with or without mapping)
> Create / Update / Delete documents
> Search documents
* The search can be done with a template to be faster than other complicated frameworks.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/elastic
```

## Usage 
This examples are available in the project at [elastic/examples](https://github.com/joaosoft/elastic/tree/master/examples)

### Templates
#### get.example.1.template
```
{
  "query": {
    "bool": {
      "must": {
        "term": {
          {{ range $key, $value := .Data }}
             "{{ $key }}": "{{ $value }}"
             {{ if (gt (len $.Data) 1) }}
                 ,
             {{ end }}
          {{ end }}
        }
      }
    }
  },
  "sort": [
    {
      "age": {
        "order": "desc"
      }
    }
  ]

  {{ if (gt $.From 0) }}
    ,
    "from": {{.From}}
  {{ end }}

  {{ if (gt $.Size 0) }}
    ,
  " size": {{.Size}}
  {{ end }}
}
```

### Code
```go
// create a client
import "github.com/joaosoft/elastic"

client := elastic.NewClient("http://localhost:9200")
// you can define the configuration without having a configuration file
//client1 := elastic.NewElastic(elastic.WithConfiguration(elastic.NewConfig("http://localhost:9200")))

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// create a index with mapping
err := client.CreateIndex().Index("persons").Body([]byte(`
{
  "mappings": {
    "person": {
      "properties": {
        "age": {
          "type": "long"
        },
        "name": {
          "type": "text",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 256
            }
          }
        }
      }
    }
  }
}
`)).Execute()


// create a new document with id
age, _ := strconv.Atoi(id)
id, err := client.Create().Index("persons").Type("person").Id(id).Body(person{
    Name: "joao",
    Age:  age + 20,
}).Execute()

if err != nil {
    log.Error(err)
} else {
    fmt.Printf("\ncreated a new person with id %s\n", id)
}

	
// create a new document with a generated id
id, err := client.Create().Index("persons").Type("person").Body(person{
    Name: "joao",
    Age:  30,
}).Execute()

if err != nil {
    log.Error(err)
} else {
    fmt.Printf("\ncreated a new person with id %s\n", id)
}


// update a document
age, _ := strconv.Atoi(id)
id, err := client.Create().Index("persons").Type("person").Id(id).Body(person{
    Name: "luis",
    Age:  age + 20,
}).Execute()

if err != nil {
    log.Error(err)
} else {
    fmt.Printf("\nupdated person with id %s\n", id)
}


// search a document with a template
d1 := elastic.TemplateData{Data: map[string]interface{}{"name": name}}

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


// delete a document
err := client.Delete().Index("persons").Type("person").Id(id).Execute()

if err != nil {
    log.Error(err)
} else {
    fmt.Printf("\ndeleted person with id %s\n", "1")
}


// validate if exists a index
status, err := client.ExistsIndex().Index(index).Execute()

if err != nil {
    log.Error(err)
} else {
    fmt.Printf("\nexists index? %t\n", status == http.StatusOK)
}


// delete a index
err := client.DeleteIndex().Index("persons").Execute()

if err != nil {
    log.Error(err)
} else {
    fmt.Printf("\ndeleted persons index\n")
} 
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
