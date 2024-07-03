package test

import (
	"fmt"
	"testing"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/client"
	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/service"
)

func setup_session() chan *client.WebClient {
	ch := make(chan *client.WebClient)

	config.Config = config.NewConfig()

	routes := service.Handler()
	cl := client.NewClient()

	srv := service.NewServer(routes)

	go func() {
		ch <- cl
		if err := srv.Server.ListenAndServe(); err != nil {
			panic(fmt.Sprintf("could not start server: %v", err))
		}
	}()

	return ch
}

func TestSession(t *testing.T) {
	ch := setup_session()
	cl := <-ch
	s, err := cl.Session(capabilities.DefaultCapabilities())
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if s == nil {
		fmt.Println("no session id")
		t.Fail()
	}

	fmt.Println(s)
}
