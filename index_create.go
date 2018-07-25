package elastic

import (
	"encoding/json"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joaosoft/errors"
)

type CreateIndexResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

type CreateIndex struct {
	client *Elastic
	index  string
	body   []byte
	method string
}

func NewCreateIndex(e *Elastic) *CreateIndex {
	return &CreateIndex{
		client: e,
		method: http.MethodPut,
	}
}

func (e *CreateIndex) Index(index string) *CreateIndex {
	e.index = index
	return e
}

func (e *CreateIndex) Body(body interface{}) *CreateIndex {
	switch v := body.(type) {
	case []byte:
		e.body = v
	default:
		e.body, _ = json.Marshal(v)
	}
	return e
}

func (e *CreateIndex) Execute() error {

	// create data on elastic
	reader := bytes.NewReader(e.body)

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s", e.client.config.Endpoint, e.index), reader)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)

	elasticResponse := CreateIndexResponse{}
	json.Unmarshal(body, &elasticResponse)

	if !elasticResponse.Acknowledged {
		return errors.New("couldn't create the index")
	}

	return nil
}
