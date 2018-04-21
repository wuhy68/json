package main

import (
	"fmt"

	"go-dropbox/service"

	"github.com/labstack/gommon/log"
)

func main() {
	access := "Bearer"
	token := "<HERE>"
	api := "https://api.dropboxapi.com/2"

	config := godropbox.NewConfig(access, token, api)
	dropbox := godropbox.NewDropbox(godropbox.WithConfiguration(config))

	// get user information
	log.Info("get user information")
	if user, err := dropbox.User().GetUserAccount(); err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("%+v", user)
	}

	// upload a file
	log.Info("upload a file")
	if err := dropbox.File().Upload("teste.txt", []byte("teste")); err != nil {
		log.Error(err.Error())
	}

	// download the uploaded file
	log.Info("download the uploaded file")
	if file, err := dropbox.File().Download("teste.jpg"); err != nil {
		log.Error(err.Error())
	} else {
		fmt.Println(file)
	}

}
