package dropbox

import logger "github.com/joaosoft/logger"

// DropboxOption ...
type DropboxOption func(dropbox *Dropbox)

// Reconfigure ...
func (dropbox *Dropbox) Reconfigure(options ...DropboxOption) {
	for _, option := range options {
		option(dropbox)
	}
}

// WithConfiguration ...
func WithConfiguration(config *DropboxConfig) DropboxOption {
	return func(dropbox *Dropbox) {
		dropbox.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILog) DropboxOption {
	return func(dropbox *Dropbox) {
		log = logger
		dropbox.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) DropboxOption {
	return func(dropbox *Dropbox) {
		log.SetLevel(level)
	}
}
