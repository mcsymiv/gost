package driver

import (
	"fmt"
	"os"
	"os/exec"
	"time"

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

	webclient := client.NewClient()

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
		panic(fmt.Sprintf("error on quit: %v", err))
	}
}

func (w *WebDriver) FindElement(selector *data.Selector) *WebElement {
	eId, err := w.WebClient.FindElement(selector, w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on find element: %v", err))
	}

	return &WebElement{
		WebDriver:    w,
		WebElementId: eId,
	}
}

func (w *WebDriver) FindElements(selector *data.Selector) []*WebElement {
	elementsId, err := w.WebClient.FindElements(selector, w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on find elements: %v", err))
	}

	var els []*WebElement

	for _, id := range elementsId {
		els = append(els, &WebElement{
			WebDriver:          w,
			WebElementId:       id,
			WebElementSelector: selector,
		})
	}

	return els
}

func (w *WebDriver) Fs(s string) []*WebElement {
	selector := Strategy(s)
	elementsId, err := w.WebClient.FindElements(selector, w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on find element: %v", err))
	}

	var els []*WebElement

	for _, id := range elementsId {
		els = append(els, &WebElement{
			WebDriver:          w,
			WebElementId:       id,
			WebElementSelector: selector,
		})
	}

	return els
}

func (w *WebDriver) F(s string) *WebElement {
	selector := Strategy(s)
	eId, err := w.WebClient.FindElement(selector, w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on find element: %v", err))
	}

	return &WebElement{
		WebDriver:          w,
		WebElementId:       eId,
		WebElementSelector: selector,
	}
}

// Next
// finds xpath element from element
func (w *WebElement) Next(s string) *WebElement {
	by := NextStrategy(s)

	eId, err := w.WebClient.FromElement(by, w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on find element: %v", err))
	}

	return &WebElement{
		WebDriver:          w.WebDriver,
		WebElementId:       eId,
		WebElementSelector: by,
	}
}

func (w *WebElement) Nexts(s string) []*WebElement {
	by := NextStrategy(s)

	elementsId, err := w.WebClient.FromElements(by, w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on find element: %v", err))
	}

	var els []*WebElement

	for _, id := range elementsId {
		els = append(els, &WebElement{
			WebDriver:          w.WebDriver,
			WebElementId:       id,
			WebElementSelector: by,
		})
	}

	return els
}

// P
// invokes parent function
// with N-level times
func (w *WebElement) Up(level int) *WebElement {
	by := PXpathStrategy(level, w.WebElementSelector)

	eId, err := w.WebClient.FromElement(by, w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on find element: %v", err))
	}

	return &WebElement{
		WebDriver:          w.WebDriver,
		WebElementId:       eId,
		WebElementSelector: by,
	}
}

func (w *WebElement) Parent() *WebElement {
	by := ParentXpathStrategy(w.WebElementSelector)

	eId, err := w.WebClient.FromElement(by, w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on find element: %v", err))
	}

	return &WebElement{
		WebDriver:          w.WebDriver,
		WebElementId:       eId,
		WebElementSelector: by,
	}
}

func (w *WebElement) IsDisplayed() bool {
	ok, err := w.WebClient.IsDisplayed(w.SessionId, w.WebElementId)
	fmt.Println(ok)
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

// Cl
// finds and clicks on element
func (w *WebDriver) Cl(s string) *WebElement {
	selector := Strategy(s)
	eId, err := w.WebClient.FindElement(selector, w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on find element: %v", err))
	}

	el := &WebElement{
		WebDriver:    w,
		WebElementId: eId,
	}

	err = w.WebClient.Click(w.SessionId, el.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on click: %v", err))
	}

	return el
}

// Input
// inputs keys, text to a input element
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

// Text
// retrieves text from element
func (w *WebElement) Text() string {
	txt, err := w.WebClient.Text(w.SessionId, w.WebElementId)
	if err != nil {
		panic(fmt.Sprintf("error on action: %v", err))
	}

	return txt
}

func (w *WebDriver) Until(fn func() bool) {
	var success bool
	start := time.Now()
	end := start.Add(config.Config.WaitForTimeout * time.Second)

	for {
		success = fn()
		if success {
			break
		}

		if time.Now().After(end) {
			panic("error on until")
		}

		time.Sleep(config.Config.WaitForInterval * time.Millisecond)
	}
}

func readFile(name string) ([]byte, error) {
	c, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("error on read file: %v", err)
	}

	return c, nil
}

// SetValueJs
// Combines selenium selector strategy
// And Find element method with JS set value
func (w *WebDriver) SetValueJs(selector, value string) {
	el := w.F(selector)

	args := []interface{}{el.WebElementId, value}

	f, err := config.FindFile(config.Config.JsFilesPath, "setValue.js")
	if err != nil {
		panic(fmt.Sprintf("error on find file: %v", err))
	}

	c, err := readFile(f)
	if err != nil {
		panic(fmt.Sprintf("error on file read in setValue.js: %v", err))
	}

	err = w.WebClient.Script(string(c), w.SessionId, args)
	if err != nil {
		panic(fmt.Sprintf("error on script: %v", err))
	}
}

// ClickJs
// Combines selenium selector strategy
// And Find element method with JS click
func (w *WebDriver) ClickJs(selector string) {
	el := w.F(selector)

	args := []interface{}{el.WebElementId}

	f, err := config.FindFile(config.Config.JsFilesPath, "click.js")
	if err != nil {
		panic(fmt.Sprintf("error on find file: %v", err))
	}

	c, err := readFile(f)
	if err != nil {
		panic(fmt.Sprintf("error on file read in setValue.js: %v", err))
	}

	err = w.WebClient.Script(string(c), w.SessionId, args)
	if err != nil {
		panic(fmt.Sprintf("error on script: %v", err))
	}
}
