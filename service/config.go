package godropbox

// appConfig ...
type appConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Authorization struct {
		access string `json:"access"`
		token  string `json:"token"`
	} `json:"authorization"`
}
