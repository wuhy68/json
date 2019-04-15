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
