package elastic

import (
	"fmt"

	gomanager "github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Elastic ElasticConfig `json:"elastic"`
}

// ElasticConfig ...
type ElasticConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Endpoint string `json:"endpoint"`
}

// NewConfig ...
func NewConfig() *ElasticConfig {
	appConfig := &AppConfig{}
	if _, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/models.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	}

	return &appConfig.Elastic
}
