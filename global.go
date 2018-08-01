package mailer

import logger "github.com/joaosoft/logger"

var log = logger.NewLogDefault("mailer", logger.InfoLevel)
var templates = make(map[string][]byte)
