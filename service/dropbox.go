package godropbox

import (
	"fmt"

	"github.com/joaosoft/go-error/service"
	"github.com/joaosoft/go-manager/service"
)

type dropbox struct {
	client gomanager.IGateway
	config appConfig

	user   *user
	file   *file
}

func NewDropbox(options ...goDropboxOption) *dropbox {
	dropbox := &dropbox{
		client: gomanager.NewSimpleGateway(),
	}
	goerror.NewError(fmt.Errorf(""))

	dropbox.reconfigure(options...)

	return dropbox
}

func (d *dropbox) User() *user {
	if d.user == nil {
		d.user = &user{dropbox: d}
	}
	return d.user
}

func (d *dropbox) File() *file {
	if d.file == nil {
		d.file = &file{dropbox: d}
	}
	return d.file
}
