package main

import (
	"fmt"
	"reflect"
	"validator"

	"errors"

	"regexp"

	"github.com/satori/go.uuid"
)

const (
	RegexForMissingParm = `%\+?[a-z]`
)

type Data string

type Example struct {
	Name       string         `validate:"value=joao, dummy, error={1:a;b}, max=10"`
	Age        int            `validate:"value=30, error={99}"`
	Street     int            `validate:"max=10, error=3"`
	Brothers   []Example      `validate:"size=1, error=4"`
	Id         uuid.UUID      `validate:"nonzero, error=5"`
	Option1    string         `validate:"options=aa;bb;cc, error=6"`
	Option2    int            `validate:"options=11;22;33, error=7"`
	Option3    []string       `validate:"options=aa;bb;cc, error=8"`
	Option4    []int          `validate:"options=11;22;33, error=9"`
	Map1       map[string]int `validate:"options=aa:11;bb:22;cc:33, error=10"`
	Map2       map[int]string `validate:"options=11:aa;22:bb;33:cc, error=11"`
	StartTime  string         `validate:"special={time}, error=12"`
	StartDate1 string         `validate:"special={date}, error=13"`
	StartDate2 string         `validate:"special={YYYYMMDD}, error=14"`
	DateString *string        `validate:"special={YYYYMMDD}, error=15"`
	Data       *Data          `validate:"special={YYYYMMDD}, error=16"`
	unexported string
}

var dummy_middle_handler = func(name string, value reflect.Value, expected interface{}, errs *[]error) []error {
	var rtnErrs []error

	err := errors.New("dummy responding...")
	rtnErrs = append(rtnErrs, err)

	return rtnErrs
}

func init() {
	validator.
		AddMiddle("dummy", dummy_middle_handler).
		SetValidateAll(true).
		SetErrorCodeHandler(dummy_error_handler)
}

var errs = map[string]error{
	"1":  errors.New("error 1: a:%s, b:%s"),
	"2":  errors.New("error 2"),
	"3":  errors.New("error 3"),
	"4":  errors.New("error 4"),
	"5":  errors.New("error 5"),
	"6":  errors.New("error 6"),
	"7":  errors.New("error 7"),
	"8":  errors.New("error 8"),
	"9":  errors.New("error 9"),
	"10": errors.New("error 10"),
	"11": errors.New("error 11"),
	"12": errors.New("error 12"),
	"13": errors.New("error 13"),
	"14": errors.New("error 14"),
	"15": errors.New("error 15"),
	"16": errors.New("error 16"),
}
var dummy_error_handler = func(code string, arguments []interface{}, name string, value reflect.Value, expected interface{}, err *[]error) error {
	if err, ok := errs[code]; ok {
		var regx = regexp.MustCompile(RegexForMissingParm)
		matches := regx.FindAllStringIndex(err.Error(), -1)

		if len(matches) > 0 {

			if len(arguments) < len(matches) {
				arguments = append(arguments, name)
			}

			err = fmt.Errorf(err.Error(), arguments...)
		}

		return err
	}
	return nil
}

func main() {
	id, _ := uuid.NewV4()
	str := "2018-12-1"
	data := Data("2018-12-1")
	example := Example{
		Id:         id,
		Name:       "joao",
		Age:        30,
		Street:     10,
		Option1:    "aa",
		Option2:    11,
		Option3:    []string{"aa", "bb", "cc"},
		Option4:    []int{11, 22, 33},
		Map1:       map[string]int{"aa": 11, "bb": 22, "cc": 33},
		Map2:       map[int]string{11: "aa", 22: "bb", 33: "cc"},
		StartTime:  "12:01:00",
		StartDate1: "01-12-2018",
		StartDate2: "2018-12-1",
		DateString: &str,
		Data:       &data,
		Brothers: []Example{
			Example{
				Name:       "jessica",
				Age:        10,
				Street:     12,
				Option1:    "xx",
				Option2:    99,
				Option3:    []string{"aa", "zz", "cc"},
				Option4:    []int{11, 44, 33},
				Map1:       map[string]int{"aa": 11, "kk": 22, "cc": 33},
				Map2:       map[int]string{11: "aa", 22: "bb", 99: "cc"},
				StartTime:  "99:01:00",
				StartDate1: "01-99-2018",
				StartDate2: "2018-99-1",
			},
		},
	}
	if errs := validator.Validate(example); len(errs) > 0 {
		fmt.Printf("ERRORS: %d\n", len(errs))
		for _, err := range errs {
			fmt.Printf("\nERROR: %s", err)
		}
	}
}
