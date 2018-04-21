package godropbox

import (
	"fmt"

	"github.com/joaosoft/go-log/service"
	"github.com/joaosoft/go-manager/service"
)

type dropbox struct {
	client gomanager.IGateway
	config *config
	pm     *gomanager.GoManager

	// usage..
	user *user
	file *file
}

// NewDropbox ...
func NewDropbox(options ...goDropboxOption) *dropbox {
	pm := gomanager.NewManager(gomanager.WithLogger(log), gomanager.WithRunInBackground(false))

	dropbox := &dropbox{
		client: gomanager.NewSimpleGateway(),
		pm:     pm,
	}

	// load configuration file
	appConfig := &appConfig{}
	if simpleConfig, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := golog.ParseLevel(appConfig.Log.Level)
		log.Debugf("setting log level to %s", level)
		WithLogLevel(level)
	}

	dropbox.reconfigure(options...)

	return dropbox
}

// User ...
func (d *dropbox) User() *user {
	if d.user == nil {
		d.user = &user{
			client: d.client,
			config: d.config,
		}
	}
	return d.user
}

// File ...
func (d *dropbox) File() *file {
	if d.file == nil {
		d.file = &file{
			client: d.client,
			config: d.config,
		}
	}
	return d.file
}
