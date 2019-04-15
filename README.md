json
================

[![Build Status](https://travis-ci.org/joaosoft/json.svg?branch=master)](https://travis-ci.org/joaosoft/json) | [![codecov](https://codecov.io/gh/joaosoft/json/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/json) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/json)](https://goreportcard.com/report/github.com/joaosoft/json) | [![GoDoc](https://godoc.org/github.com/joaosoft/json?status.svg)](https://godoc.org/github.com/joaosoft/json)

A simple json marshal and unmarshal by customized tags (exported fields only).

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for


## Dependecy Management
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/json
```

## Usage 
This examples are available in the project at [json/examples](https://github.com/joaosoft/json/tree/master/examples)

### Code
```go
package main

import (
	"fmt"
	"json"
)

type address struct {
	Street string            `db.read:"street"`
	Number float64           `db.write:"number"`
	Map    map[string]string `db:"map"`
}

type person struct {
	Name    string              `db:"name"`
	Age     int                 `db:"age"`
	Address *address            `db:"address"`
	Numbers []int               `db:"numbers"`
	Others  map[string]*address `db:"others"`
}

func main() {
	marshal()
	unmarshal()
}

func marshal() {
	fmt.Println("\n\n:: MARSHAL")

	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:    "joao",
		Age:     30,
		Address: addr,
		Numbers: []int{1, 2, 3},
		Others:  map[string]*address{`"ola" "joao"`: addr, "c": addr},
	}

	// with tags "db" and "db.read"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))

	// with tags "db" and "db.write"
	// marshal
	bytes, err = json.Marshal(example, "db", "db.write")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func unmarshal() {
	fmt.Println("\n\n:: UNMARSHAL")

	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:    "joao",
		Age:     30,
		Address: addr,
		Numbers: []int{1, 2, 3},
		Others:  map[string]*address{`"ola" "joao"`: addr, "c": addr},
	}

	// with tags "db" and "db.read"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))

	// unmarshal
	var newExample person
	err = json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n:: Example: %+v", newExample)
	fmt.Printf("\n:: Address: %+v\n\n\n", newExample.Address)

	// with tags "db" and "db.write"
	// marshal
	bytes, err = json.Marshal(example, "db", "db.write")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))

	// unmarshal
	newExample = person{}
	err = json.Unmarshal(bytes, &newExample, "db", "db.write")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: Example: %+v", newExample)
	fmt.Printf("\n:: Address: %+v", newExample.Address)

	for key, value := range newExample.Others {
		fmt.Printf("\n:: Others Key: %+v", key)
		fmt.Printf("\n:: Others Value: %+v", value)
	}
}
```

> ##### Result:
```go
:: MARSHAL
{"name":"joao","age":30,"address":{"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"numbers":[1,2,3],"others":{"\"ola\" \"joao\"":{"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"c":{"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}}}
{"name":"joao","age":30,"address":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"numbers":[1,2,3],"others":{"\"ola\" \"joao\"":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"c":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}}}


:: UNMARSHAL
{"name":"joao","age":30,"address":{"street":"street one","map":{"c":"d","\"ola\" \"joao\"":"\"adeus\" \"joao\""}},"numbers":[1,2,3],"others":{"\"ola\" \"joao\"":{"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"c":{"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}}}

:: Example: {Name:joao Age:30 Address:0xc00000a1c0 Numbers:[1 2 3] Others:map["ola" "joao":0xc00000a280 c:0xc00000a2a0]}
:: Address: &{Street:street one Number:0 Map:map[c:d "ola" "joao":"adeus" "joao"]}


{"name":"joao","age":30,"address":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"numbers":[1,2,3],"others":{"c":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"\"ola\" \"joao\"":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}}}

:: Example: {Name:joao Age:30 Address:0xc00000a360 Numbers:[1 2 3] Others:map[c:0xc00000a420 "ola" "joao":0xc00000a440]}
:: Address: &{Street: Number:1.2 Map:map["ola" "joao":"adeus" "joao" c:d]}
:: Others Key: c
:: Others Value: &{Street: Number:1.2 Map:map["ola" "joao":"adeus" "joao" c:d]}
:: Others Key: "ola" "joao"
:: Others Value: &{Street: Number:1.2 Map:map["ola" "joao":"adeus" "joao" c:d]}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
