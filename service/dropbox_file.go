package godropbox

import (
	"github.com/joaosoft/go-error/service"
	"github.com/joaosoft/go-manager/service"
)

type file struct {
	client gomanager.IGateway
	config *config
}

func (f *file) Download(file string) ([]byte, *goerror.ErrorData) {
	return nil, nil
}

func (f *file) Upload(name string, content []byte) *goerror.ErrorData {
	return nil
}
