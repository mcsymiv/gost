package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/command"
	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/service"
)

func startDriver(capsFn ...capabilities.CapabilitiesFunc) (*driver.WebDriver, func()) {
	d := driver.Driver()
	return d, func() {
		// teardown
		d.Quit()
		command.OutFileLogs.Close()
		d.Command.Process.Kill()
	}
}

func setup_quit() func() {
	config.Config = config.NewConfig()

	routes := service.Handler()
	srv := service.NewServer(routes)

	go func() {
		if err := srv.Server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			panic(fmt.Sprintf("could not start server: %v", err))
		}
	}()

	return func() {
		shutdownCtx := context.Background()

		if err := srv.Server.Shutdown(shutdownCtx); err != nil {
			panic(fmt.Sprintf("HTTP shutdown error: %v", err))
		}
		fmt.Println("Graceful shutdown complete")
	}
}

func TestQuit(t *testing.T) {
	shutdown := setup_quit()
	defer shutdown()

	wd, tear := startDriver()
	defer tear()

	wd.Open("https://google.com")
}
