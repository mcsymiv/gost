package test

import (
	"fmt"
	"testing"

	"github.com/mcsymiv/gost/client"
	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/service"
)

func setup_status() chan *client.WebClient {
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

func TestStatus(t *testing.T) {
	ch := setup_status()
	cl := <-ch
	st, _ := cl.Status()

	fmt.Println(st)
}
