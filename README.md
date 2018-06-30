# dropbox
[![Build Status](https://travis-ci.org/joaosoft/dropbox.svg?branch=master)](https://travis-ci.org/joaosoft/dropbox) | [![codecov](https://codecov.io/gh/joaosoft/dropbox/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/dropbox) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/dropbox)](https://goreportcard.com/report/github.com/joaosoft/dropbox) | [![GoDoc](https://godoc.org/github.com/joaosoft/dropbox?status.svg)](https://godoc.org/github.com/joaosoft/dropbox)

A simple dropbox v2 client.

## Support for 
> User
* Get account information

> Files
* Upload / Download files
* Create / Delete files

>Folders
* List files
* Create folders
* Delete folders

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/dropbox
```

## Usage 
This examples are available in the project at [dropbox/examples](https://github.com/joaosoft/dropbox/tree/master/examples)
```go
import "github.com/joaosoft/dropbox"

dropbox := dropbox.NewDropbox()

//get user information
log.Info("get user information")
if user, err := dropbox.User().Get(); err != nil {
    log.Error(err.Error())
} else {
    fmt.Printf("\n\nUSER: %+v \n\n", user)
}

// upload a file
log.Info("upload a file")
if response, err := dropbox.File().Upload("/teste.txt", []byte("teste")); err != nil {
    log.Error(err.Error())
} else {
    fmt.Printf("\n\nUPLOADED: %+v \n\n", response)
}

// download the uploaded file
log.Info("download the uploaded file")
if response, err := dropbox.File().Download("/teste.txt"); err != nil {
    log.Error(err.Error())
} else {
    fmt.Printf("\n\nDOWNLOADED: %s \n\n", string(response))
}

// create folder
log.Info("listing folder")
if response, err := dropbox.Folder().Create("/bananas"); err != nil {
    log.Error(err.Error())
} else {
    fmt.Printf("\n\nCREATED FOLDER: %+v \n\n", response)
}

// listing folder
log.Info("listing folder")
if response, err := dropbox.Folder().List("/"); err != nil {
    log.Error(err.Error())
} else {
    fmt.Printf("\n\nLIST FOLDER: %+v \n\n", response)
}

// deleting the uploaded file
log.Info("deleting the uploaded file")
if response, err := dropbox.File().Delete("/teste.txt"); err != nil {
    log.Error(err.Error())
} else {
    fmt.Printf("\n\nDELETED FILE: %+v \n\n", response)
}

// deleting the created folder
log.Info("deleting the created folder")
if response, err := dropbox.Folder().DeleteFolder("/bananas"); err != nil {
    log.Error(err.Error())
} else {
    fmt.Printf("\n\nDELETED FOLDER: %+v \n\n", response)
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
