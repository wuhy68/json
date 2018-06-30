package dropbox

import (
	"fmt"

	gomanager "github.com/joaosoft/manager")

// AppConfig ...
type AppConfig struct {
	Dropbox DropboxConfig `json:"dropbox"`
}

// DropboxConfig ...
type DropboxConfig struct {
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
func NewConfig(access, token string) *DropboxConfig {
	appConfig := &AppConfig{}
	if _, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/models.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	}

	appConfig.Dropbox.Authorization.Access = access
	appConfig.Dropbox.Authorization.Token = token

	return &appConfig.Dropbox
}
