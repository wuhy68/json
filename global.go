package dropbox

import logger "github.com/joaosoft/logger"

var global = make(map[string]interface{})
var log = logger.NewLogDefault("dropbox", logger.InfoLevel)

func init() {
	global[path_key] = defaultPath
}
