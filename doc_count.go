package elastic

import (
	"encoding/json"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	errors "github.com/joaosoft/errors"
)

type CountResponse struct {
	count int64 `json:"count"`
}

type Count struct {
	client *Elastic
	index  string
	typ    string
	query  string
	method string
}

func NewCount(client *Elastic) *Count {
	return &Count{
		client: client,
		method: http.MethodGet,
	}
}

func (e *Count) Index(index string) *Count {
	e.index = index
	return e
}

func (e *Count) Type(typ string) *Count {
	e.typ = typ
	return e
}

func (e *Count) Query(query string) *Count {
	e.query = query
	return e
}

type CountTemplate struct {
	Data interface{} `json:"data,omitempty"`
}

func (e *Count) Template(path, name string, data *CountTemplate, reload bool) *Count {
	key := fmt.Sprintf("%s/%s", path, name)

	var result bytes.Buffer
	var err error

	if _, found := templates[key]; !found {
		e.client.mux.Lock()
		defer e.client.mux.Unlock()
		templates[key], err = ReadFile(key, nil)
		if err != nil {
			log.Error(err)
			return e
		}
	}

	t := template.New(name)
	t, err = t.Parse(string(templates[key]))
	if err == nil {
		if err := t.ExecuteTemplate(&result, name, data); err != nil {
			log.Error(err)
			return e
		}

		e.query = result.String()
	} else {
		log.Error(err)
		return e
	}

	return e
}

func (e *Count) Execute() (int64, error) {

	if e.query != "" {
		e.method = http.MethodPost
	}

	// get data from elastic
	reader := strings.NewReader(e.query)

	var q string
	if e.typ != "" {
		q += fmt.Sprintf("/%s", e.typ)
	}

	q += "/_count"

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s%s", e.client.config.Endpoint, e.index, q), reader)
	if err != nil {
		return 0, errors.New(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error(err)
		return 0, errors.New(err)
	}
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return 0, errors.New(err)
	}

	elasticResponse := CountResponse{}
	if err := json.Unmarshal(body, &elasticResponse); err != nil {
		log.Error(err)
		return 0, errors.New(err)
	}

	return elasticResponse.count, nil
}
