package mailer

import (
	logger "github.com/joaosoft/logger"
)

// MailerOption ...
type MailerOption func(client *Mailer)

// Reconfigure ...
func (mailer *Mailer) Reconfigure(options ...MailerOption) {
	for _, option := range options {
		option(mailer)
	}
}

// WithConfiguration ...
func WithConfiguration(config *MailerConfig) MailerOption {
	return func(client *Mailer) {
		client.config = config
		client.auth = PlainAuth(client.config.Identity, client.config.Username, client.config.Password, client.config.Host)
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) MailerOption {
	return func(mailer *Mailer) {
		log = logger
		mailer.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) MailerOption {
	return func(mailer *Mailer) {
		log.SetLevel(level)
	}
}
