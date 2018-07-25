package elastic

import (
	"encoding/json"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joaosoft/errors"
)

type CreateScrollResponse struct {
	Index   string `json:"_index,omitempty"`
	Type    string `json:"_type,omitempty"`
	ID      string `json:"_id,omitempty"`
	Version int64  `json:"_version,omitempty"`
	Found   bool   `json:"found"`
	Result  string `json:"result"`
	Created bool   `json:"created"`
}

type CreateScroll struct {
	client *Elastic
	index  string
	typ    string
	body   []byte
	method string
}

func NewCreateScroll(e *Elastic) *CreateScroll {
	return &CreateScroll{
		client: e,
		method: http.MethodPut,
	}
}

func (e *CreateScroll) Index(index string) *CreateScroll {
	e.index = index
	return e
}

func (e *CreateScroll) Type(typ string) *CreateScroll {
	e.typ = typ
	return e
}

func (e *CreateScroll) Body(body interface{}) *CreateScroll {
	switch v := body.(type) {
	case []byte:
		e.body = v
	default:
		e.body, _ = json.Marshal(v)
	}
	return e
}

func (e *CreateScroll) Execute() (string, error) {

	// create data on elastic
	reader := bytes.NewReader(e.body)

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s/%s", e.client.config.Endpoint, e.index, e.typ), reader)
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)

	elasticResponse := CreateScrollResponse{}
	json.Unmarshal(body, &elasticResponse)

	if !elasticResponse.Created && elasticResponse.Result != "updated" {
		return "", errors.New("couldn't create the scroll")
	}

	return elasticResponse.ID, nil
}
