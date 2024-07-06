package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/service"
)

func setup_url() func() {
	config.Config = config.NewConfig()

	routes := service.Handler()
	srv := service.NewServer(routes)

	go func() {
		if err := srv.Server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			panic(fmt.Sprintf("could not start server: %v", err))
		}
	}()

	return func() {
		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 2*time.Second)
		defer shutdownRelease()

		if err := srv.Server.Shutdown(shutdownCtx); err != nil {
			panic(fmt.Sprintf("HTTP shutdown error: %v", err))
		}
		fmt.Println("Graceful shutdown complete")
	}
}

func TestUrl(t *testing.T) {
	shutdown := setup_url()
	defer shutdown()

	wd := driver.Driver()
	wd.Open("https://google.com")
	time.Sleep(5 * time.Second)
}
