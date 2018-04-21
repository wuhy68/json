package godropbox

// appConfig ...
type appConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}

type config struct {
	Authorization struct {
		Access string `json:"access"`
		Token  string `json:"token"`
	} `json:"authorization"`
	Api string `json:"api"`
}

// NewConfig ...
func NewConfig(access, token, api string) *config {
	config := &config{Api: api}
	config.Authorization.Access = access
	config.Authorization.Token = token

	return config
}
