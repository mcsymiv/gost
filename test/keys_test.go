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

func start_keys(capsFn ...capabilities.CapabilitiesFunc) (*driver.WebDriver, func()) {
	d := driver.Driver()
	// setup
	return d, func() {
		// teardown
		d.Quit()
		command.OutFileLogs.Close()
		d.Command.Process.Kill()
	}
}

func setup_keys() func() {
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

func TestKeys(t *testing.T) {
	shutdown := setup_keys()
	defer shutdown()

	wd, tear := start_keys(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	wd.Open("https://google.com")
	el := wd.F("//*[@id='APjFqb']").Click()

	if el == nil {
		t.Fail()
	}

	clickedEl := el.Input("hello")
	if clickedEl == nil {
		t.Fail()
	}

	if el.WebElementId != clickedEl.WebElementId {
		t.Fail()
	}

	fmt.Println(clickedEl.WebElementId)
}
