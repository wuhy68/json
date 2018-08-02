package main

import (
	"fmt"
	"mailer"
	"os"
	"time"
)

var client = mailer.NewMailer()

func main() {

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf(err.Error())
	}

	image, err := mailer.ReadFile(dir+"/examples/attachments/mail.png", nil)
	failed, err := client.SendMessage().
		From("João Ribeiro", "joaosoft@gmail.com").
		To("joao.ribeiro@foursource.pt", "invalid", "joao.ribeiro@foursource.pt").
		Cc("joao.ribeiro@foursource.pt", "joao.ribeiro@foursource.pt").
		Bcc("joao.ribeiro@foursource.pt", "joao.ribeiro@foursource.pt").
		Header("aFrom", "fake@mail.pt").
		Subject("This is a test subject").
		Body(mailer.ContentTypeTextPlain, "Hello, you got an email!\n\n").
		Date(time.Now()).
		Attachment(image, true, "image_file_1.png").
		Execute()

	if err != nil {
		fmt.Printf(err.Error())
	}

	if len(failed) > 0 {
		fmt.Printf("\n\nFailed addresses: %+v", failed)
	}
}
