package godropbox

import "github.com/joaosoft/go-error/service"

type file struct {
	*dropbox
}

func (f *file) Download(file string) ([]byte, *goerror.ErrorData) {
	return nil, nil
}

func (f *file) Upload(name string, content []byte) *goerror.ErrorData {
	return nil
}
