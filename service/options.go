package godropbox

import "github.com/joaosoft/go-log/service"

// dropboxOption ...
type dropboxOption func(dropbox *Dropbox)

// Reconfigure ...
func (dropbox *Dropbox) Reconfigure(options ...dropboxOption) {
	for _, option := range options {
		option(dropbox)
	}
}

// WithConfiguration ...
func WithConfiguration(config *DropboxConfig) dropboxOption {
	return func(dropbox *Dropbox) {
		dropbox.config = config
	}
}

// WithLogger ...
func WithLogger(logger golog.ILog) dropboxOption {
	return func(dropbox *Dropbox) {
		log = logger
		dropbox.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level golog.Level) dropboxOption {
	return func(dropbox *Dropbox) {
		log.SetLevel(level)
	}
}
