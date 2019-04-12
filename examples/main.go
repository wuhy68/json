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
	Name    string   `db:"name"`
	Age     int      `db:"age"`
	Address *address `db:"address"`
}

func main() {
	marshal()
	unmarshal()
}

func marshal() {
	fmt.Println("\n\n:: MARSHAL")

	example := person{Name: "joao", Age: 30, Address: &address{Street: "street one", Number: 1.2, Map: map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"}}}

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
	fmt.Println(":: UNMARSHAL")

	example := person{Name: "joao", Age: 30, Address: &address{Street: "street one", Number: 1.2}}

	// with tags "db" and "db.read"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}

	// unmarshal
	var newExample person
	err = json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", newExample)

	// with tags "db" and "db.write"
	// marshal
	bytes, err = json.Marshal(example, "db", "db.write")
	if err != nil {
		panic(err)
	}

	// unmarshal
	newExample = person{}
	err = json.Unmarshal(bytes, &newExample, "db", "db.write")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", newExample)
}
