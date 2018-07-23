package elastic

import logger "github.com/joaosoft/logger"

var log = logger.NewLogDefault("elastic", logger.InfoLevel)
var templates = make(map[string][]byte)
