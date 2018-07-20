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

type ElasticResponse struct {
	Hits   Hits   `json:"hits,omitempty"`
	Error  *Error `json:"error"`
	Status int    `json:"status"`
}

type Hits struct {
	Total int   `json:"total,omitempty"`
	Hits  []Hit `json:"hits,omitempty"`
}

type Hit struct {
	Index   string          `json:"_index,omitempty"`
	Type    string          `json:"_type,omitempty"`
	ID      string          `json:"_id,omitempty"`
	Version int64           `json:"_version,omitempty"`
	Found   bool            `json:"found"`
	Source  json.RawMessage `json:"_source,omitempty"`
}

type Error struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type Search struct {
	client   *Elastic
	index    string
	document string
	query    string
	object   interface{}
	method   string
}

func NewSearch(client *Elastic) *Search {
	return &Search{
		client: client,
		method: http.MethodGet,
	}
}

func (elastic *Elastic) Search() *Search {
	return NewSearch(elastic)
}

func (e *Search) Index(index string) *Search {
	e.index = index
	return e
}

func (e *Search) Document(document string) *Search {
	e.document = document
	return e
}

func (e *Search) Query(query string) *Search {
	if query != "" {
		e.method = http.MethodPost
	}
	e.query = query
	return e
}

func (e *Search) Object(object interface{}) *Search {
	e.object = object
	return e
}

type TemplateData struct {
	Data interface{} `json:"data,omitempty"`
	From int         `json:"from,omitempty"`
	Size int         `json:"size,omitempty"`
}

func (e *Search) Template(path, name string, data *TemplateData, reload bool) *Search {
	key := fmt.Sprintf("%s/%s", path, name)

	var result bytes.Buffer
	if _, ok := templates[key]; !ok {
		templates[key], _ = readFile(key, nil)
	}

	t := template.New(name)
	t, err := t.Parse(string(templates[key]))
	if err == nil {

		t.ExecuteTemplate(&result, name, data)

		e.query = result.String()
	}

	return e
}

func (e *Search) Execute() error {

	// get data from elastic
	reader := strings.NewReader(e.query)
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/%s/_search", e.client.config.Endpoint, e.index, e.document), reader)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)
	elasticResponse := ElasticResponse{}
	json.Unmarshal(body, &elasticResponse)

	if elasticResponse.Error != nil {
		return errors.FromString(fmt.Sprintf("[%s] %s", elasticResponse.Error.Type, elasticResponse.Error.Reason))
	}

	rawHits := make([]json.RawMessage, len(elasticResponse.Hits.Hits))
	for i, hit := range elasticResponse.Hits.Hits {
		rawHits[i] = hit.Source
	}

	hit, err := json.Marshal(rawHits)
	if err != nil {
		return errors.NewError(err)
	}

	if err := json.Unmarshal(hit, e.object); err != nil {
		return errors.NewError(err)
	}

	return nil
}
