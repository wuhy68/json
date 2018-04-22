package godropbox

import (
	"fmt"

	"github.com/joaosoft/go-log/service"
	"github.com/joaosoft/go-manager/service"
)

type Dropbox struct {
	client gomanager.IGateway
	config *GoDropboxConfig
	pm     *gomanager.GoManager

	// usage ...
	user   *user
	folder *folder
	file   *file
}

// NewDropbox ...
func NewDropbox(options ...goDropboxOption) *Dropbox {
	pm := gomanager.NewManager(gomanager.WithRunInBackground(false))

	// load configuration file
	appConfig := &appConfig{}
	if simpleConfig, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := golog.ParseLevel(appConfig.GoDropbox.Log.Level)
		log.Debugf("setting log level to %s", level)
		WithLogLevel(level)
	}

	dropbox := &Dropbox{
		client: gomanager.NewSimpleGateway(),
		pm:     pm,
		config: &appConfig.GoDropbox,
	}

	dropbox.reconfigure(options...)

	return dropbox
}

// Api ...
func (d *Dropbox) User() *user {
	if d.user == nil {
		d.user = &user{
			client: d.client,
			config: d.config,
		}
	}
	return d.user
}

// Folder ...
func (d *Dropbox) Folder() *folder {
	if d.folder == nil {
		d.folder = &folder{
			client: d.client,
			config: d.config,
		}
	}
	return d.folder
}

// File ...
func (d *Dropbox) File() *file {
	if d.file == nil {
		d.file = &file{
			client: d.client,
			config: d.config,
		}
	}
	return d.file
}
