# mailer
[![Build Status](https://travis-ci.org/joaosoft/mailer.svg?branch=master)](https://travis-ci.org/joaosoft/mailer) | [![codecov](https://codecov.io/gh/joaosoft/mailer/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/mailer) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/mailer)](https://goreportcard.com/report/github.com/joaosoft/mailer) | [![GoDoc](https://godoc.org/github.com/joaosoft/mailer?status.svg)](https://godoc.org/github.com/joaosoft/mailer)

A simple and fast SMTP email client.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/mailer
```

## Usage 
This examples are available in the project at [mailer/examples](https://github.com/joaosoft/mailer/tree/master/examples)

### Code
```go
import "github.com/joaosoft/mailer"

// create a client
client := mailer.NewMailer()

dir, _ := os.Getwd()

image, err := mailer.ReadFile(dir+"/examples/attachments/mail.png", nil)
failed, err := client.SendMessage().
    From("João Ribeiro", "joaosoft@gmail.com").
    To("joao.ribeiro@foursource.pt", "invalid", "joao.ribeiro@foursource.pt").
    Cc("joao.ribeiro@foursource.pt", "joao.ribeiro@foursource.pt").
    Bcc("joao.ribeiro@foursource.pt", "joao.ribeiro@foursource.pt").
    Header("aFrom", "fake@mail.pt").
    Subject("This is a test subject").
    Body("Hello, you got an email!\n\n").
    Date(time.Now()).
    Attachment(image, true, "image_file_1.png").
    Execute()

if err != nil {
    fmt.Printf(err.Error())
}

if len(failed) > 0 {
    fmt.Printf("\n\nFailed addresses: %+v", failed)
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
