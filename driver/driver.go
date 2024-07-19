package driver

import (
	"fmt"
	"os/exec"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/client"
	"github.com/mcsymiv/gost/command"
	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/data"
)

type Element interface{}

type WebDriver struct {
	Command      *exec.Cmd
	WebClient    *client.WebClient
	Capabilities *capabilities.Capabilities
	SessionId    string
}

type WebElement struct {
	*WebDriver
	WebElementId       string
	WebElementSelector *data.Selector
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
	url, err := w.WebClient.Url(u, w.SessionId)
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

func (w *WebDriver) Quit() {
	err := w.WebClient.Quit(w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on open: %v", err))
	}
}

func (w *WebDriver) FindElement(selector *data.Selector) *WebElement {
	eId, err := w.WebClient.FindElement(selector, w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on open: %v", err))
	}

	return &WebElement{
		WebElementId: eId,
	}
}

func (w *WebDriver) F(s string) *WebElement {
	selector := Strategy(s)
	eId, err := w.WebClient.FindElement(selector, w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on find element: %v", err))
	}

	return &WebElement{
		WebDriver:    w,
		WebElementId: eId,
	}
}

func (w *WebElement) IsDisplayed() bool {
	ok, err := w.WebClient.IsDisplayed(w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on isdisplayed: %v", err))
	}

	return ok
}

func (w *WebElement) Is() *WebElement {
	ok, err := w.WebClient.Is(w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on is: %v", err))
	}

	if !ok {
		panic(fmt.Sprintf("error on is, element not displayed: %v", err))
	}

	return w
}

func (w *WebElement) Click() *WebElement {
	err := w.WebClient.Click(w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on click: %v", err))
	}

	return w
}

func (w *WebElement) Keys(keys string) *WebElement {
	err := w.WebClient.Keys(keys, w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on click: %v", err))
	}

	return w
}
func (w *WebElement) Attr(attr string) string {
	a, err := w.WebClient.Attr(attr, w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on click: %v", err))
	}

	return a
}
