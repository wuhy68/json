package godropbox

import "github.com/joaosoft/go-log/service"

// goDropboxOption ...
type goDropboxOption func(godropbox *Dropbox)

// reconfigure ...
func (godropbox *Dropbox) reconfigure(options ...goDropboxOption) {
	for _, option := range options {
		option(godropbox)
	}
}

// WithConfiguration ...
func WithConfiguration(config *goDropboxConfig) goDropboxOption {
	return func(godropbox *Dropbox) {
		godropbox.config = config
	}
}

// WithLogLevel ...
func WithLogLevel(level golog.Level) goDropboxOption {
	return func(godropbox *Dropbox) {
		log.SetLevel(level)
	}
}
