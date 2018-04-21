package godropbox

import "github.com/joaosoft/go-log/service"

// goDropboxOption ...
type goDropboxOption func(godropbox *dropbox)

// reconfigure ...
func (godropbox *dropbox) reconfigure(options ...goDropboxOption) {
	for _, option := range options {
		option(godropbox)
	}
}

// WithConfiguration ...
func WithConfiguration(config *config) goDropboxOption {
	return func(godropbox *dropbox) {
		godropbox.config = config
	}
}

// WithLogLevel ...
func WithLogLevel(level golog.Level) goDropboxOption {
	return func(godropbox *dropbox) {
		log.SetLevel(level)
	}
}
