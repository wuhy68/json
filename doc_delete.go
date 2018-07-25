package elastic

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joaosoft/errors"
)

type DeleteResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

type DeleteHit struct {
	Found  bool   `json:"found"`
	Result string `json:"result"`
}

type Delete struct {
	client *Elastic
	index  string
	typ    string
	id     string
	method string
}

func NewDelete(client *Elastic) *Delete {
	return &Delete{
		client: client,
		method: http.MethodDelete,
	}
}

func (e *Delete) Index(index string) *Delete {
	e.index = index
	return e
}

func (e *Delete) Type(typ string) *Delete {
	e.typ = typ
	return e
}

func (e *Delete) Id(id string) *Delete {
	e.id = id
	return e
}

func (e *Delete) Execute() error {

	// delete data from elastic
	var query string
	if e.typ != "" {
		query += fmt.Sprintf("/%s", e.typ)
	}

	if e.id != "" {
		query += fmt.Sprintf("/%s", e.id)
	}

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s%s", e.client.config.Endpoint, e.index, query), nil)
	if err != nil {
		return errors.New(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New(err)
	}
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New(err)
	}

	if e.id != "" {
		elasticResponse := DeleteHit{}
		if err = json.Unmarshal(body, &elasticResponse); err != nil {
			return errors.New(err)
		}

		if !elasticResponse.Found || elasticResponse.Result != "deleted" {
			return errors.New("couldn't delete the resource")
		}
	} else {
		elasticResponse := DeleteResponse{}
		if err = json.Unmarshal(body, &elasticResponse); err != nil {
			return errors.New(err)
		}

		if !elasticResponse.Acknowledged {
			return errors.New("couldn't delete the resource")
		}
	}

	return nil
}
