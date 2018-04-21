package godropbox

import "github.com/joaosoft/go-error/service"

type user struct {
	*dropbox
}

type userData struct {
	Name string
}

func (u *user) GetUser() (*userData, *goerror.ErrorData) {
	return &userData{Name: "joao"}, nil
}
