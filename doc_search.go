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

type SearchResponse struct {
	Hits   SearchHits   `json:"hits,omitempty"`
	Error  *SearchError `json:"error"`
	Status int          `json:"status"`
}

type SearchHits struct {
	Total int         `json:"total,omitempty"`
	Hits  []SearchHit `json:"hits,omitempty"`
}

type SearchHit struct {
	Index   string          `json:"_index,omitempty"`
	Type    string          `json:"_type,omitempty"`
	ID      string          `json:"_id,omitempty"`
	Version int64           `json:"_version,omitempty"`
	Found   bool            `json:"found"`
	Source  json.RawMessage `json:"_source,omitempty"`
}

type SearchError struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type Search struct {
	client *Elastic
	index  string
	typ    string
	id     string
	query  string
	object interface{}
	method string
}

func NewSearch(client *Elastic) *Search {
	return &Search{
		client: client,
		method: http.MethodGet,
	}
}

func (e *Search) Index(index string) *Search {
	e.index = index
	return e
}

func (e *Search) Type(typ string) *Search {
	e.typ = typ
	return e
}

func (e *Search) Id(id string) *Search {
	e.id = id
	return e
}

func (e *Search) Query(query string) *Search {
	e.query = query
	return e
}

func (e *Search) Object(object interface{}) *Search {
	e.object = object
	return e
}

type SearchTemplate struct {
	Data interface{} `json:"data,omitempty"`
	From int         `json:"from,omitempty"`
	Size int         `json:"size,omitempty"`
}

func (e *Search) Template(path, name string, data *SearchTemplate, reload bool) *Search {
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

func (e *Search) Execute() error {

	if e.query != "" {
		e.method = http.MethodPost
	}

	// get data from elastic
	reader := strings.NewReader(e.query)

	var q string
	if e.typ != "" {
		q += fmt.Sprintf("/%s", e.typ)
	}

	if e.id != "" {
		q += fmt.Sprintf("/%s", e.id)
	} else {
		q += "/_search"
	}

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s%s", e.client.config.Endpoint, e.index, q), reader)
	if err != nil {
		return errors.New(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error(err)
		return errors.New(err)
	}
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return errors.New(err)
	}

	var hit []byte

	if e.id != "" {
		elasticResponse := SearchHit{}
		if err := json.Unmarshal(body, &elasticResponse); err != nil {
			log.Error(err)
			return errors.New(err)
		}

		hit, err = json.Marshal(elasticResponse.Source)
		if err != nil {
			log.Error(err)
			return errors.New(err)
		}
	} else {
		elasticResponse := SearchResponse{}
		if err := json.Unmarshal(body, &elasticResponse); err != nil {
			log.Error(err)
			return errors.New(err)
		}

		if elasticResponse.Error != nil {
			return errors.New(fmt.Sprintf("[%s] %s", elasticResponse.Error.Type, elasticResponse.Error.Reason))
		}

		rawHits := make([]json.RawMessage, len(elasticResponse.Hits.Hits))
		for i, rawHit := range elasticResponse.Hits.Hits {
			rawHits[i] = rawHit.Source
		}

		hit, err = json.Marshal(rawHits)
		if err != nil {
			return errors.New(err)
		}
	}

	if err := json.Unmarshal(hit, e.object); err != nil {
		return errors.New(err)
	}

	return nil
}
