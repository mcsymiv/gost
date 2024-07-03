package test

import (
	"fmt"
	"testing"

	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/service"
)

func TestService(t *testing.T) {
	config.Config = config.NewConfig()
	routes := service.Handler()
	srv := service.NewServer(routes)

	if err := srv.Server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("could not start server: %v", err))
	}
}
