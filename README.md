# elastic
[![Build Status](https://travis-ci.org/joaosoft/elastic.svg?branch=master)](https://travis-ci.org/joaosoft/elastic) | [![codecov](https://codecov.io/gh/joaosoft/elastic/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/elastic) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/elastic)](https://goreportcard.com/report/github.com/joaosoft/elastic) | [![GoDoc](https://godoc.org/github.com/joaosoft/elastic?status.svg)](https://godoc.org/github.com/joaosoft/elastic)

A simple elastic client.

## Support for 
> Search
* Search with templates

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
  }

  {{ if (gt $.Size 0) }}
  ,
  {{ end }}
  {{ if (gt $.From 0) }}
  "from": {{.From}}
  {{ end }}
  {{ if (gt $.Size 0) }}
  ,
  {{ end }}
  {{ if (gt $.Size 0) }}
  "size": {{.Size}}
  {{ end }}
}
```

#### get.example.2.template
```
{
  "query": {
    "bool": {
      "filter": {
        "match": {
          "ids": {
            "query": "{{ .Data }}"
          }
        }
      }
    }
  }

  {{ if (gt $.Size 0) }}
  ,
  {{ end }}
  {{ if (gt $.From 0) }}
  "from": {{.From}}
  {{ end }}
  {{ if (gt $.Size 0) }}
  ,
  {{ end }}
  {{ if (gt $.Size 0) }}
  "size": {{.Size}}
  {{ end }}
}
```

### Code
```go
import "github.com/joaosoft/elastic"

var data []interface{}

client := elastic.NewClient("http://localhost:9200")
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
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
