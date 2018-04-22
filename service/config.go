package godropbox

import (
	"fmt"

	"github.com/joaosoft/go-manager/service"
)

// appConfig ...
type appConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Authorization struct {
		Access string `json:"access"`
		Token  string `json:"token"`
	} `json:"authorization"`
	Hosts struct {
		Api     string `json:"api"`
		Content string `json:"content"`
	} `json:"hosts"`
}

// NewConfig ...
func NewConfig(access, token string) *appConfig {
	appConfig := &appConfig{}
	if _, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	}

	appConfig.Authorization.Access = access
	appConfig.Authorization.Token = token

	return appConfig
}
