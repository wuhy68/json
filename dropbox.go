package dropbox

import (
	"fmt"

	logger "github.com/joaosoft/logger"
	manager "github.com/joaosoft/manager"
)

type Dropbox struct {
	client        manager.IGateway
	config        *DropboxConfig
	pm            *manager.Manager
	isLogExternal bool

	// usage ...
	user   *User
	folder *Folder
	file   *File
}

// NewDropbox ...
func NewDropbox(options ...DropboxOption) *Dropbox {
	pm := manager.NewManager(manager.WithRunInBackground(false))

	dropbox := &Dropbox{
		client: manager.NewSimpleGateway(),
		pm:     pm,
	}

	dropbox.Reconfigure(options...)

	if dropbox.isLogExternal {
		pm.Reconfigure(manager.WithLogger(log))
	}

	// load configuration File
	appConfig := &AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/models.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.Dropbox.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(logger.WithLevel(level))
	}

	dropbox.config = &appConfig.Dropbox

	return dropbox
}

// Api ...
func (d *Dropbox) User() *User {
	if d.user == nil {
		d.user = &User{
			client: d.client,
			config: d.config,
		}
	}
	return d.user
}

// Folder ...
func (d *Dropbox) Folder() *Folder {
	if d.folder == nil {
		d.folder = &Folder{
			client: d.client,
			config: d.config,
		}
	}
	return d.folder
}

// File ...
func (d *Dropbox) File() *File {
	if d.file == nil {
		d.file = &File{
			client: d.client,
			config: d.config,
		}
	}
	return d.file
}
