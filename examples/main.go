package main

import (
	j "encoding/json"
	"fmt"
	"json"
	"reflect"
	"time"
)

type address struct {
	Ports     []int             `db.read:"ports"`
	Street    string            `db.read:"street"`
	Number    float64           `db:"number" db.write:"number"`
	Timestamp time.Time         `db:"timestamp"`
	Map       map[string]string `db:"map"`
}

type person struct {
	Name      string              `db:"name"`
	Age       int                 `db:"age"`
	Address   *address            `db:"address"`
	Numbers   []int               `db:"numbers"`
	Others    map[string]*address `db:"others"`
	Addresses []*address          `db:"addresses"`
}

type contents []content

type content struct {
	Name string        `db:"name"`
	Data *j.RawMessage `db:"data"`
}

type operationType string

type operation struct {
	Operation operationType `db:"operation"`
	Name      *string       `db:"name"`
	Contacts  *contacts     `db:"contacts"`
}

type contacts struct {
	Country      string                 `db:"country"`
	Addresses    map[string]string      `db:"addresses"`
	PhoneNumbers map[string]interface{} `db:"phone_numbers"`
}

type operationList []*operation

func main() {
	marshal()
	unmarshal()
}

func marshal() {
	fmt.Println("\n\n:: MARSHAL")

	marshal_example_1()
	marshal_example_2()

}

func unmarshal() {
	fmt.Println("\n\n:: UNMARSHAL")

	unmarshal_example_1()
	unmarshal_example_2()
	unmarshal_example_3()
	unmarshal_example_4()
	unmarshal_example_5()
	unmarshal_example_6()
	unmarshal_example_7()
	unmarshal_example_8()
	unmarshal_example_9()
}

func marshal_example_1() {
	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:      "joao",
		Age:       30,
		Address:   addr,
		Numbers:   []int{1, 2, 3},
		Others:    map[string]*address{`"ola" "joao"`: addr, "c": addr},
		Addresses: []*address{addr, addr},
	}

	// with tags "db" and "db.read"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func marshal_example_2() {
	addr := &address{
		Street:    "street one",
		Number:    1.2,
		Timestamp: time.Now(),
		Map:       map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:      "joao",
		Age:       30,
		Address:   addr,
		Numbers:   []int{1, 2, 3},
		Others:    map[string]*address{`"ola" "joao"`: addr, "c": addr},
		Addresses: []*address{addr, addr},
	}

	// with tags "db" and "db.write"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.write")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func marshal_example_3() {
	data := j.RawMessage([]byte(`"data":{"test": "one", "test": "two"}`))
	example := []content{
		{
			Name: "joao",
			Data: &data,
		},
	}

	// with tags "db" and "db.read"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func unmarshal_example_1() {
	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:      "joao",
		Age:       30,
		Address:   addr,
		Numbers:   []int{1, 2, 3},
		Others:    map[string]*address{`"ola" "joao"`: addr, "c": addr},
		Addresses: []*address{addr, addr},
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
}

func unmarshal_example_2() {
	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:      "joao",
		Age:       30,
		Address:   addr,
		Numbers:   []int{1, 2, 3},
		Others:    map[string]*address{`"ola" "joao"`: addr, "c": addr},
		Addresses: []*address{addr, addr},
	}

	// with tags "db" and "db.write"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.write")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))

	// unmarshal
	newExample := person{}
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

	for _, value := range newExample.Addresses {
		fmt.Printf("\n:: Addresses: %+v", value)
	}
}

func unmarshal_example_3() {
	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:      "joao",
		Age:       30,
		Address:   addr,
		Numbers:   []int{1, 2, 3},
		Others:    map[string]*address{`"ola" "joao"`: addr, "c": addr},
		Addresses: []*address{addr, addr},
	}

	persons := []*person{&example, &example}
	bytes, err := json.Marshal(persons, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n %s", string(bytes))

	// unmarshal
	var newPersons []*person
	err = json.Unmarshal(bytes, &newPersons, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: LEN: %d", len(newPersons))
	fmt.Printf("\n:: Example 1: %+v", newPersons[0])
	fmt.Printf("\n:: Example 1 Address: %+v", newPersons[0].Address)
	fmt.Printf("\n:: Example 2: %+v", newPersons[1])
}

func unmarshal_example_4() {
	example := []int{1, 2, 3}

	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n %s", string(bytes))

	// unmarshal
	var newExample []int
	err = json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: LEN: %d", len(newExample))
	fmt.Printf("\n:: Example: %+v", newExample)
}

func unmarshal_example_5() {
	example := []int{}

	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n %s", string(bytes))

	// unmarshal
	var newExample []int
	err = json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: LEN: %d", len(newExample))
	fmt.Printf("\n:: Example: %+v", newExample)
}

func unmarshal_example_6() {
	example := map[string]int{"one": 1, "two": 2, "three": 3}

	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n %s", string(bytes))

	// unmarshal
	var newExample map[string]int
	err = json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: LEN: %d", len(newExample))
	fmt.Printf("\n:: Example: %+v", newExample)
}

func unmarshal_example_7() {
	bytes := []byte(`[{"name": "joao", "data":{"test": "one", "test": "two"}}]`)

	// unmarshal
	var newExample contents
	err := json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: LEN: %d", len(newExample))
	fmt.Printf("\n:: Example: %+v", newExample[0].Data)
}

func unmarshal_example_8() {
	bytes := []byte(`{"name":"joao","age":30,"address":{"street":"one","number":7}}`)

	// unmarshal
	var person person
	err := json.Unmarshal(bytes, &person, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: Person: %+v", person)
	fmt.Printf("\n:: Address: %+v", person.Address)
}

func unmarshal_example_9() {
	name := "joao"
	operList1 := operationList{
		&operation{
			Operation: "test",
			Name:      &name,
			Contacts: &contacts{
				Country:      "portugal",
				Addresses:    map[string]string{"casa": "1111111", "trabalho": "2222222"},
				PhoneNumbers: map[string]interface{}{"float64": 1, "boolean": false, "string": "ola", "object": contacts{Country: "test"}},
			},
		},
	}

	// marshal
	fmt.Println("MARSHAL")
	bytes, err := json.Marshal(operList1, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))

	// unmarshal
	fmt.Println("UNMARSHAL")
	var operList2 operationList
	err = json.Unmarshal(bytes, &operList2, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: Operation: %+v, Addresses: %+v, Phone Numbers: %+v, Type string:%s, Type boolean: %s, Type float64: %s",
		operList2[0].Operation,
		operList2[0].Contacts.Addresses,
		operList2[0].Contacts.PhoneNumbers,
		reflect.TypeOf(operList2[0].Contacts.PhoneNumbers["string"]),
		reflect.TypeOf(operList2[0].Contacts.PhoneNumbers["boolean"]),
		reflect.TypeOf(operList2[0].Contacts.PhoneNumbers["float64"]),
	)
}
