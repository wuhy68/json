package elastic

import (
	"encoding/json"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joaosoft/errors"
)

type CreateResponse struct {
	Index   string `json:"_index,omitempty"`
	Type    string `json:"_type,omitempty"`
	ID      string `json:"_id,omitempty"`
	Version int64  `json:"_version,omitempty"`
	Found   bool   `json:"found"`
	Result  string `json:"result"`
	Created bool   `json:"created"`
}

type Create struct {
	client *Elastic
	index  string
	typ    string
	id     string
	body   []byte
	method string
}

func NewCreate(e *Elastic) *Create {
	return &Create{
		client: e,
		method: http.MethodPost,
	}
}

func (e *Create) Index(index string) *Create {
	e.index = index
	return e
}

func (e *Create) Type(typ string) *Create {
	e.typ = typ
	return e
}

func (e *Create) Id(id string) *Create {
	e.id = id
	return e
}

func (e *Create) Body(body interface{}) *Create {
	switch v := body.(type) {
	case []byte:
		e.body = v
	default:
		e.body, _ = json.Marshal(v)
	}
	return e
}

func (e *Create) Execute() (string, error) {

	// create data on elastic
	reader := bytes.NewReader(e.body)

	var query string

	if e.id != "" {
		query += fmt.Sprintf("/%s", e.id)
	}

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s/%s%s", e.client.config.Endpoint, e.index, e.typ, query), reader)
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)

	elasticResponse := CreateResponse{}
	json.Unmarshal(body, &elasticResponse)

	if !elasticResponse.Created && elasticResponse.Result != "updated" {
		return "", errors.FromString("couldn't create the resource")
	}

	return elasticResponse.ID, nil
}
