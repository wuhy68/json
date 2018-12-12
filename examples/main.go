package main

import (
	"fmt"
	"validator"

	"errors"

	"regexp"

	"github.com/satori/go.uuid"
)

const (
	RegexForMissingParms = `%\+?[a-z]`
)

type Data string

type NextSet struct {
	Set int `validate:"set=321, id=next_set"`
}

type Example struct {
	Name               string         `validate:"value=joao, dummy_middle, error={ErrorTag1:a;b}, max=10"`
	Age                int            `validate:"value=30, error={ErrorTag99}"`
	Street             int            `validate:"max=10, error={ErrorTag3}"`
	Brothers           []Example2     `validate:"size=1, error={ErrorTag4}"`
	Id                 uuid.UUID      `validate:"nonzero, error={ErrorTag5}"`
	Option1            string         `validate:"options=aa;bb;cc, error={ErrorTag6}"`
	Option2            int            `validate:"options=11;22;33, error={ErrorTag7}"`
	Option3            []string       `validate:"options=aa;bb;cc, error={ErrorTag8}"`
	Option4            []int          `validate:"options=11;22;33, error={ErrorTag9}"`
	Map1               map[string]int `validate:"options=aa:11;bb:22;cc:33, error={ErrorTag10}"`
	Map2               map[int]string `validate:"options=11:aa;22:bb;33:cc, error={ErrorTag11}"`
	SpecialTime        string         `validate:"special=time, error={ErrorTag12}"`
	SpecialDate1       string         `validate:"special=date, error={ErrorTag13}"`
	SpecialDate2       string         `validate:"special=YYYYMMDD, error={ErrorTag14}"`
	SpecialDateString  *string        `validate:"special=YYYYMMDD, error={ErrorTag15}"`
	SpecialData        *Data          `validate:"special=YYYYMMDD, error={ErrorTag16}"`
	SpecialUrl         string         `validate:"special=url"`
	unexported         string
	IsNill             *string `validate:"nonzero, error={ErrorTag17}"`
	Sanitize           string  `validate:"sanitize=a;b;teste, error={ErrorTag17}"`
	Callback           string  `validate:"callback=dummy_callback, error={ErrorTag19}"`
	Password           string  `json:"password" validate:"id=password"`
	PasswordConfirm    string  `validate:"match=password"`
	MyName             string  `validate:"id=name"`
	MyAge              int     `validate:"id=age"`
	MyValidate         int     `validate:"if=(id=age value=30) or (id=age value=31) and (id=name value=joao), value=10"`
	DoubleValidation   int     `validate:"nonzero, error=20, min=5, error={ErrorTag21}"`
	Set                int     `validate:"set=321, id=set"`
	NextSet            NextSet
	DistinctIntPointer []*int    `validate:"distinct"`
	DistinctInt        []int     `validate:"distinct"`
	DistinctString     []string  `validate:"distinct"`
	DistinctBool       []bool    `validate:"distinct"`
	DistinctFloat      []float32 `validate:"distinct"`
	IsZero             int       `validate:"iszero"`
	Trim               string    `validate:"trim"`
	KeyValue           string    `validate:"id=my_value"`
	Key                string    `validate:"key={my_value}"`
}

type Example2 struct {
	Name              string         `validate:"value=joao, dummy_middle, error={ErrorTag1:a;b}, max=10"`
	Age               int            `validate:"value=30, error={ErrorTag99}"`
	Street            int            `validate:"max=10, error={ErrorTag3}"`
	Id                uuid.UUID      `validate:"nonzero, error={ErrorTag5}"`
	Option1           string         `validate:"options=aa;bb;cc, error={ErrorTag6}"`
	Option2           int            `validate:"options=11;22;33, error={ErrorTag7}"`
	Option3           []string       `validate:"options=aa;bb;cc, error={ErrorTag8}"`
	Option4           []int          `validate:"options=11;22;33, error={ErrorTag9}"`
	Map1              map[string]int `validate:"options=aa:11;bb:22;cc:33, error={ErrorTag10}"`
	Map2              map[int]string `validate:"options=11:aa;22:bb;33:cc, error={ErrorTag11}"`
	SpecialTime       string         `validate:"special=time, error={ErrorTag12}"`
	SpecialDate1      string         `validate:"special=date, error={ErrorTag13}"`
	SpecialDate2      string         `validate:"special=YYYYMMDD, error={ErrorTag14}"`
	SpecialDateString *string        `validate:"special=YYYYMMDD, error={ErrorTag15}"`
	SpecialData       *Data          `validate:"special=YYYYMMDD, error={ErrorTag16}"`
	SpecialUrl        string         `validate:"special=url"`
	unexported        string
	IsNill            *string `validate:"nonzero, error={ErrorTag17}"`
	Sanitize          string  `validate:"sanitize=a;b;teste, error={ErrorTag17}"`
	Callback          string  `validate:"callback=dummy_callback, error={ErrorTag19}"`
	Password          string  `json:"password" validate:"id=password"`
	PasswordConfirm   string  `validate:"match=password"`
}

var dummy_middle_handler = func(context *validator.ValidatorContext, validationData *validator.ValidationData) []error {
	var rtnErrs []error

	err := errors.New("dummy middle responding...")
	rtnErrs = append(rtnErrs, err)

	return rtnErrs
}

func init() {
	validator.
		AddMiddle("dummy_middle", dummy_middle_handler).
		SetValidateAll(true).
		SetErrorCodeHandler(dummy_error_handler).
		AddCallback("dummy_callback", dummy_callback)
}

var errs = map[string]error{
	"ErrorTag1":  errors.New("error 1: a:%s, b:%s"),
	"ErrorTag2":  errors.New("error 2"),
	"ErrorTag3":  errors.New("error 3"),
	"ErrorTag4":  errors.New("error 4"),
	"ErrorTag5":  errors.New("error 5"),
	"ErrorTag6":  errors.New("error 6"),
	"ErrorTag7":  errors.New("error 7"),
	"ErrorTag8":  errors.New("error 8"),
	"ErrorTag9":  errors.New("error 9"),
	"ErrorTag10": errors.New("error 10"),
	"ErrorTag11": errors.New("error 11"),
	"ErrorTag12": errors.New("error 12"),
	"ErrorTag13": errors.New("error 13"),
	"ErrorTag14": errors.New("error 14"),
	"ErrorTag15": errors.New("error 15"),
	"ErrorTag16": errors.New("error 16"),
	"ErrorTag17": errors.New("error 17"),
	"ErrorTag18": errors.New("error 18"),
	"ErrorTag19": errors.New("error 19"),
	"ErrorTag20": errors.New("error 20"),
	"ErrorTag21": errors.New("error 21"),
}
var dummy_error_handler = func(context *validator.ValidatorContext, validationData *validator.ValidationData) error {
	if err, ok := errs[validationData.ErrorData.Code]; ok {
		var regx = regexp.MustCompile(RegexForMissingParms)
		matches := regx.FindAllStringIndex(err.Error(), -1)

		if len(matches) > 0 {

			if len(validationData.ErrorData.Arguments) < len(matches) {
				validationData.ErrorData.Arguments = append(validationData.ErrorData.Arguments, validationData.Name)
			}

			err = fmt.Errorf(err.Error(), validationData.ErrorData.Arguments...)
		}

		return err
	}
	return nil
}

var dummy_callback = func(context *validator.ValidatorContext, validationData *validator.ValidationData) []error {
	return []error{errors.New("there's a bug here!")}
}

func main() {
	intVal1 := 1
	intVal2 := 2
	id, _ := uuid.NewV4()
	str := "2018-12-1"
	data := Data("2018-12-1")
	example := Example{
		Id:                id,
		Name:              "joao",
		Age:               30,
		Street:            10,
		Option1:           "aa",
		Option2:           11,
		Option3:           []string{"aa", "bb", "cc"},
		Option4:           []int{11, 22, 33},
		Map1:              map[string]int{"aa": 11, "bb": 22, "cc": 33},
		Map2:              map[int]string{11: "aa", 22: "bb", 33: "cc"},
		SpecialTime:       "12:01:00",
		SpecialDate1:      "01-12-2018",
		SpecialDate2:      "2018-12-1",
		SpecialDateString: &str,
		SpecialData:       &data,
		SpecialUrl:        "xxx.xxx.teste.pt",
		Password:          "password",
		PasswordConfirm:   "password_errada",
		MyName:            "joao",
		MyAge:             30,
		MyValidate:        30,
		DoubleValidation:  0,
		Set:               123,
		NextSet: NextSet{
			Set: 123,
		},
		DistinctIntPointer: []*int{&intVal1, &intVal1, &intVal2, &intVal2},
		DistinctInt:        []int{1, 1, 2, 2},
		DistinctString:     []string{"a", "a", "b", "b"},
		DistinctBool:       []bool{true, true, false, false},
		DistinctFloat:      []float32{1.1, 1.1, 1.2, 1.2},
		Trim:               "     aqui       tem     espaços    !!   ",
		KeyValue:           "     aaaaa     3245 79 / ( ) ? =  tem     espaços ...   !!  <<<< ",
		Brothers: []Example2{
			Example2{
				Name:            "jessica",
				Age:             10,
				Street:          12,
				Option1:         "xx",
				Option2:         99,
				Option3:         []string{"aa", "zz", "cc"},
				Option4:         []int{11, 44, 33},
				Map1:            map[string]int{"aa": 11, "kk": 22, "cc": 33},
				Map2:            map[int]string{11: "aa", 22: "bb", 99: "cc"},
				SpecialTime:     "99:01:00",
				SpecialDate1:    "01-99-2018",
				SpecialDate2:    "2018-99-1",
				Sanitize:        "b teste",
				SpecialUrl:      "http://www.teste.pt",
				Password:        "password",
				PasswordConfirm: "password",
			},
		},
	}

	fmt.Printf("\nBEFORE SET: %d", example.Set)
	fmt.Printf("\nBEFORE NEXT SET: %d", example.NextSet.Set)

	fmt.Printf("\nBEFORE DISTINCT INT POINTER: %+v", example.DistinctIntPointer)
	fmt.Printf("\nBEFORE DISTINCT INT: %+v", example.DistinctInt)
	fmt.Printf("\nBEFORE DISTINCT STRING: %+v", example.DistinctString)
	fmt.Printf("\nBEFORE DISTINCT BOOL: %+v", example.DistinctBool)
	fmt.Printf("\nBEFORE DISTINCT FLOAT: %+v", example.DistinctFloat)
	if errs := validator.Validate(&example); len(errs) > 0 {
		fmt.Printf("\n\nERRORS: %d\n", len(errs))
		for _, err := range errs {
			fmt.Printf("\nERROR: %s", err)
		}
	}
	fmt.Printf("\n\nAFTER SET: %d", example.Set)
	fmt.Printf("\nAFTER NEXT SET: %d", example.NextSet.Set)

	fmt.Printf("\nAFTER DISTINCT INT POINTER: %+v", example.DistinctIntPointer)
	fmt.Printf("\nAFTER DISTINCT INT: %+v", example.DistinctInt)
	fmt.Printf("\nAFTER DISTINCT STRING: %+v", example.DistinctString)
	fmt.Printf("\nAFTER DISTINCT BOOL: %+v", example.DistinctBool)
	fmt.Printf("\nAFTER DISTINCT FLOAT: %+v", example.DistinctFloat)

	fmt.Printf("\nAFTER TRIM: %s", example.Trim)
	fmt.Printf("\nAFTER KEY: %s", example.Key)
}
