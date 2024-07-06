package driver

import (
	"fmt"
	"os/exec"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/client"
	"github.com/mcsymiv/gost/command"
	"github.com/mcsymiv/gost/config"
)

type WebDriver struct {
	Command      *exec.Cmd
	WebClient    *client.WebClient
	Capabilities *capabilities.Capabilities
	SessionId    string
}

func NewDriver(capsFn ...capabilities.CapabilitiesFunc) *WebDriver {
	caps := capabilities.DefaultCapabilities()
	for _, capFn := range capsFn {
		capFn(caps)
	}

	webclient := client.NewClient()
	session, err := webclient.Session(caps)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &WebDriver{
		Capabilities: caps,
		WebClient:    webclient,
		SessionId:    session.Id,
	}
}

func (w *WebDriver) DriverSession() *WebDriver {
	session, err := w.WebClient.Session(w.Capabilities)
	if err != nil {
		return nil
	}

	w.SessionId = session.Id
	return w
}

func Driver(capsFn ...capabilities.CapabilitiesFunc) *WebDriver {
	caps := capabilities.DefaultCapabilities()
	for _, capFn := range capsFn {
		capFn(caps)
	}

	webclient := client.WebDriverClient()

	// TODO: start driver
	exec, err := command.Cmd(caps, config.Config)
	if err != nil {
		panic(fmt.Sprintf("error on starting driver command: %v", err))
	}

	session, err := webclient.Session(caps)
	if err != nil {
		panic(fmt.Sprintf("error on session create: %v", err))
	}

	return &WebDriver{
		Command:      exec,
		Capabilities: caps,
		WebClient:    webclient,
		SessionId:    session.Id,
	}
}

func (w *WebDriver) Url(u string) string {
	url, err := w.WebClient.Url(u)
	if err != nil {
		fmt.Printf("error on new driver session: %v", err)
		return ""
	}

	return url.Url
}

func (w *WebDriver) Open(u string) string {
	url, err := w.WebClient.Open(u, w.SessionId)
	if err != nil {
		fmt.Printf("error on open: %v", err)
		return ""
	}

	return url.Url
}

func (w *WebDriver) FindElement(selector string) *WebElement {

	return &WebElement{}
}
