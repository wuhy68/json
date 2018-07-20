package elastic

import (
	"fmt"

	logger "github.com/joaosoft/logger"
	manager "github.com/joaosoft/manager"
)

type Elastic struct {
	config        *ElasticConfig
	isLogExternal bool
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
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.Elastic.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(logger.WithLevel(level))
	}

	elastic.config = &appConfig.Elastic

	elastic.Reconfigure(options...)

	return elastic
}
