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

func start_click(capsFn ...capabilities.CapabilitiesFunc) (*driver.WebDriver, func()) {
	d := driver.Driver()
	// setup
	return d, func() {
		// teardown
		d.Quit()
		command.OutFileLogs.Close()
		d.Command.Process.Kill()
	}
}

func setup_click() func() {
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

func TestClick(t *testing.T) {
	shutdown := setup_click()
	defer shutdown()

	wd, tear := start_click(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	wd.Open("https://google.com")
	el := wd.F("//*[@id='APjFqb']")

	if el == nil {
		t.Fail()
	}

	fmt.Println(el.WebElementId)

	clickedEl := el.Click()
	if clickedEl == nil {
		t.Fail()
	}

	// if el.WebElementId != clickedEl.WebElementId {
	// 	t.Fail()
	// }

	fmt.Println(clickedEl.WebElementId)
}
