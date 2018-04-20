# go-error
[![Build Status](https://travis-ci.org/joaosoft/go-error.svg?branch=master)](https://travis-ci.org/joaosoft/go-error) | [![codecov](https://codecov.io/gh/joaosoft/go-error/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/go-error) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/go-error)](https://goreportcard.com/report/github.com/joaosoft/go-error) | [![GoDoc](https://godoc.org/github.com/joaosoft/go-error?status.svg)](https://godoc.org/github.com/joaosoft/go-error/service)

Error manager with error and caused by structure.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/go-error/service
```

## Usage 
This examples are available in the project at [go-error/example](https://github.com/joaosoft/go-error/tree/master/example)
```go
err := goerror.NewError(fmt.Errorf("erro 1"))
err.Add(fmt.Errorf("erro 2"))
err.Add(fmt.Errorf("erro 3"))

fmt.Printf("Error: %s, Cause: %s", err.Error(), err.Cause())
```

##### Result:
```javascript
Error: erro 3, Cause: 'erro 3', caused by 'erro 2', caused by 'erro 2'
```

## Known issues


## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
