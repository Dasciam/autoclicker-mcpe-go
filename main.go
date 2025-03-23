package main

import (
	"github.com/dasciam/autoclicker-mcpe-go/application"
)

func main() {
	app, err := application.New()
	if err != nil {
		panic(err)
	}
	if err := app.Run(); err != nil {
		return
	}
}
