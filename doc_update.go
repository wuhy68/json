package elastic

import (
	"encoding/json"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joaosoft/errors"
)

type UpdateResponse struct {
	Index   string `json:"_index,omitempty"`
	Type    string `json:"_type,omitempty"`
	ID      string `json:"_id,omitempty"`
	Version int64  `json:"_version,omitempty"`
	Found   bool   `json:"found"`
	Result  string `json:"result"`
	Created bool   `json:"created"`
}

type Update struct {
	client *Elastic
	index  string
	typ    string
	id     string
	body   []byte
	method string
}

func NewUpdate(e *Elastic) *Update {
	return &Update{
		client: e,
		method: http.MethodPut,
	}
}

func (e *Update) Index(index string) *Update {
	e.index = index
	return e
}

func (e *Update) Type(typ string) *Update {
	e.typ = typ
	return e
}

func (e *Update) Id(id string) *Update {
	e.id = id
	return e
}

func (e *Update) Body(body interface{}) *Update {
	e.body, _ = json.Marshal(body)
	return e
}

func (e *Update) Execute() (string, error) {

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

	elasticResponse := UpdateResponse{}
	json.Unmarshal(body, &elasticResponse)

	if elasticResponse.Result != "created" && elasticResponse.Result != "updated" {
		return "", errors.FromString("couldn't update the resource")
	}

	return elasticResponse.ID, nil
}
