package mailer

import (
	"fmt"
	"sync"

	logger "github.com/joaosoft/logger"
	manager "github.com/joaosoft/manager"
)

type Mailer struct {
	config        *MailerConfig
	auth          Auth
	isLogExternal bool
	mux           sync.Mutex
}

// NewMailer ...
func NewMailer(options ...MailerOption) *Mailer {
	pm := manager.NewManager(manager.WithRunInBackground(false))

	mailer := &Mailer{}

	if mailer.isLogExternal {
		pm.Reconfigure(manager.WithLogger(log))
	}

	// load configuration File
	appConfig := &AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.Mailer.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(logger.WithLevel(level))
	}

	mailer.config = &appConfig.Mailer
	mailer.auth = PlainAuth(mailer.config.Identity, mailer.config.Username, mailer.config.Password, mailer.config.Host)

	mailer.Reconfigure(options...)

	return mailer
}

func (e *Mailer) SendMessage() *SendMessageService {
	return NewSendMessageService(e)
}
