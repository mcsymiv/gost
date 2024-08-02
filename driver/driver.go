package driver

import (
	"fmt"
	"os"
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
		fmt.Printf("error on url: %v", err)
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

func (w *WebDriver) NewTab() {
	err := w.WebClient.NewTab(w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on new tabs: %v", err))
	}
}

func (w *WebDriver) Tabs() []string {
	tabs, err := w.WebClient.Tabs(w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on tabs: %v", err))
	}

	return tabs
}

func (w *WebDriver) Tab(n int) {
	err := w.WebClient.Tab(n, w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on tab: %v", err))
	}
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
		WebDriver:    w,
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

func (w *WebElement) Input(keys string) *WebElement {
	err := w.WebClient.Input(keys, w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on keys: %v", err))
	}

	return w
}

func (w *WebElement) Attr(attr string) string {
	a, err := w.WebClient.Attr(attr, w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on attribute: %v", err))
	}

	return a
}

func (w *WebDriver) ExecuteScript(s string, args ...interface{}) {
	err := w.WebClient.Script(s, w.SessionId, args)
	if err != nil {
		panic(fmt.Sprintf("error on script: %v", err))
	}
}

func (w *WebDriver) Script(fName string, args ...interface{}) {
	f, err := config.FindFile(config.Config.JsFilesPath, fmt.Sprintf("%s.js", fName))
	if err != nil {
		panic(fmt.Sprintf("error on find file: %v", err))
	}

	c, err := os.ReadFile(f)
	if err != nil {
		panic(fmt.Sprintf("error on read file: %v", err))
	}

	err = w.WebClient.Script(string(c), w.SessionId, args)
	if err != nil {
		panic(fmt.Sprintf("error on script: %v", err))
	}
}

func (w *WebElement) Id() map[string]string {
	return map[string]string{
		config.WebElementIdentifier: w.WebElementId,
	}
}

func (w *WebDriver) Screenshot() {
	err := w.WebClient.Screenshot(w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on script: %v", err))
	}
}

func (w *WebDriver) Active() *WebElement {
	eId, err := w.WebClient.Active(w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on find element: %v", err))
	}

	return &WebElement{
		WebDriver:    w,
		WebElementId: eId,
	}
}

func (w *WebElement) Text() string {
	txt, err := w.WebClient.Text(w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on action: %v", err))
	}

	return txt
}
