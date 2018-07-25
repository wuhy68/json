package elastic

import (
	"encoding/json"

	"bytes"
	"fmt"
	"net/http"

	"github.com/joaosoft/errors"
)

type ExistsIndex struct {
	client *Elastic
	index  string
	body   []byte
	method string
}

func NewExistsIndex(e *Elastic) *ExistsIndex {
	return &ExistsIndex{
		client: e,
		method: http.MethodHead,
	}
}

func (e *ExistsIndex) Index(index string) *ExistsIndex {
	e.index = index
	return e
}

func (e *ExistsIndex) Body(body interface{}) *ExistsIndex {
	switch v := body.(type) {
	case []byte:
		e.body = v
	default:
		e.body, _ = json.Marshal(v)
	}
	return e
}

func (e *ExistsIndex) Execute() (int, error) {

	// create data on elastic
	reader := bytes.NewReader(e.body)

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s", e.client.config.Endpoint, e.index), reader)
	if err != nil {
		return 0, err
	}

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	if err != nil {
		return response.StatusCode, errors.New(err)
	}

	return response.StatusCode, nil
}
