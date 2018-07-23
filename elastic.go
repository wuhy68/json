package elastic

import (
	"fmt"

	"sync"

	logger "github.com/joaosoft/logger"
	manager "github.com/joaosoft/manager"
)

type Elastic struct {
	config        *ElasticConfig
	isLogExternal bool
	mux           sync.Mutex
}

// NewElastic ...
func NewElastic(options ...ElasticOption) *Elastic {
	pm := manager.NewManager(manager.WithRunInBackground(false))

	elastic := &Elastic{}

	if elastic.isLogExternal {
		pm.Reconfigure(manager.WithLogger(log))
	}

	// load configuration File
	appConfig := &AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.Elastic.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(logger.WithLevel(level))
	}

	elastic.config = &appConfig.Elastic

	elastic.Reconfigure(options...)

	if elastic.config.Endpoint == "" {
		elastic.config.Endpoint = DefaultURL
	}

	return elastic
}

func (e *Elastic) Create() *Create {
	return NewCreate(e)
}

func (e *Elastic) Update() *Update {
	return NewUpdate(e)
}

func (elastic *Elastic) Delete() *Delete {
	return NewDelete(elastic)
}

func (e *Elastic) Search() *Search {
	return NewSearch(e)
}

func (e *Elastic) ExistsIndex() *ExistsIndex {
	return NewExistsIndex(e)
}

func (e *Elastic) CreateIndex() *CreateIndex {
	return NewCreateIndex(e)
}

func (e *Elastic) UpdateIndex() *CreateIndex {
	return NewCreateIndex(e)
}

func (e *Elastic) DeleteIndex() *DeleteIndex {
	return NewDeleteIndex(e)
}
