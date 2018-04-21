package main

import (
	"fmt"
	"go-dropbox/service"

	"github.com/labstack/gommon/log"
)

func main() {

	dropbox := godropbox.NewDropbox()

	// get user information
	log.Info("get user information")
	if user, err := dropbox.User().GetUser(); err != nil {
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
