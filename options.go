package elastic

import (
	logger "github.com/joaosoft/logger"
)

// ElasticOption ...
type ElasticOption func(client *Elastic)

// Reconfigure ...
func (elastic *Elastic) Reconfigure(options ...ElasticOption) {
	for _, option := range options {
		option(elastic)
	}
}

// WithConfiguration ...
func WithConfiguration(config *ElasticConfig) ElasticOption {
	return func(client *Elastic) {
		client.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) ElasticOption {
	return func(elastic *Elastic) {
		log = logger
		elastic.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) ElasticOption {
	return func(elastic *Elastic) {
		log.SetLevel(level)
	}
}
