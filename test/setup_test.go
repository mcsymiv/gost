package test

import (
	"fmt"

	"github.com/mcsymiv/gost/client"
	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/service"
)

func setup() chan *client.WebClient {
	ch := make(chan *client.WebClient)
	config.Config = config.NewConfig()

	routes := service.Handler()
	srv := service.NewServer(routes)
	cl := client.NewClient()

	go func() {
		ch <- cl
		if err := srv.Server.ListenAndServe(); err != nil {
			panic(fmt.Sprintf("could not start server: %v", err))
		}
	}()

	return ch
}

func setup2() chan *driver.WebDriver {
	ch := make(chan *driver.WebDriver)
	config.Config = config.NewConfig()

	routes := service.Handler()
	srv := service.NewServer(routes)
	cl := client.NewClient()
	cl.Status()
	wd := driver.Driver()

	go func() {
		ch <- wd
		if err := srv.Server.ListenAndServe(); err != nil {
			panic(fmt.Sprintf("could not start server: %v", err))
		}
	}()

	return ch
}
