package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/data"
)

var (
	ErrorElementId        = "error on map element id.\nValue: %v.\nError: %v"
	ErrorFindElement      = "error on find element id.\n Value: %v.\nError: %v"
	ErrorClick            = "error on click element.\nError: %v"
	ErrorSendKeys         = "error on send keys.\nError: %v"
	ErrorAttribute        = "error on attribute element.\nError: %v"
	ErrorScriptExecute    = "error on script execute.\nError: %v"
	ErrorScreenshot       = "error on screenshot.\nError: %v"
	ErrorActiveElement    = "error on active element.\nValue: %v.\nError: %v"
	ErrorTextElement      = "error on text element.\nError: %v"
	ErrorAction           = "error on action.\nError: %v"
	ErrorDisplayedElement = "error on displayed element.\nError: %v"
	ErrorDeleteSession    = "error on delete session.\nError: %v"
	ErrorCreateSession    = "error on create session.\nError: %v"
	ErrorStatus           = "error on webdriver status.\nError: %v"
	ErrorTab              = "error on tabs.\nError: %v"
	ErrorOpenUrl          = "error on open url.\nError: %v"
)

const (
	// W3C Session
	statusEndpoint     = "%s/status"
	sessionEndpoint    = "%s/session"
	quitEndpoint       = "%s/session/%s"
	urlEndpoint        = "%s/session/%s/url"
	screenshotEndpoint = "%s/session/%s/screenshot"

	// W3C Element
	findElementEndpoint  = "%s/session/%s/element"
	findElementsEndpoint = "%s/session/%s/elements"
	activeEndpoint       = "%s/session/%s/element/active"
	textEndpoint         = "%s/session/%s/element/%s/text"
	isDisplayedEndpoint  = "%s/session/%s/element/%s/displayed"
	clickEndpoint        = "%s/session/%s/element/%s/click"
	sendKeysEndpoint     = "%s/session/%s/element/%s/value"
	attributeEndpoint    = "%s/session/%s/element/%s/attribute/%s"
	fromElementEndpoint  = "%s/session/%s/element/%s/element"
	fromElementsEndpoint  = "%s/session/%s/element/%s/elements"

	// W3C Window
	windowEndpoint        = "%s/session/%s/window"
	newWindowEndpoint     = "%s/session/%s/window/new"
	windowHandlesEndpoint = "%s/session/%s/window/handles"

	// W3C Action
	actionEndpoint = "%s/session/%s/actions"

	// GoST
	isEndpoint         = "%s/session/%s/element/%s/is"
	syncScriptEndpoint = "%s/session/%s/script"
)

// RestClient represents a REST client configuration.
type WebClient struct {
	Close bool
	Host  string

	// the client UserAgent string
	UserAgent string

	// Common headers to be passed on each request
	Headers map[string]string

	// Cookies to be passed on each request
	Cookies []*http.Cookie

	// if FollowRedirects is false, a 30x response will be returned as is
	FollowRedirects bool

	// if HeadRedirects is true, the client will follow the redirect also for HEAD requests
	HeadRedirects bool

	// if Verbose, log request and response info
	Verbose bool

	WebConfig          *config.WebConfig
	WebServerAddr      string
	HTTPClient         *http.Client
	RequestReaderLimit int64
	// syncMutex  sync.Mutex // Mutex for ensuring thread safety
}

// newClientV2
// new client init without Session param
func NewClient() *WebClient {
	return &WebClient{
		WebConfig:          config.Config,
		RequestReaderLimit: 4096,

		HTTPClient: &http.Client{
			Transport: &retry{
				maxRetries: 3,
				delay:      time.Duration(config.Config.WaitForInterval * time.Millisecond),
				next: &loggin{
					next: http.DefaultTransport,
				},
			},
		},
	}
}

func marshalData(body interface{}) []byte {
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error on marshal: ", err)
		return nil
	}

	return b
}

func unmarshalRes(res *http.Response, any interface{}) error {
	return json.NewDecoder(res.Body).Decode(any)
}

type HttpResponse struct {
	http.Response
}

func ElementID(v map[string]string) (string, error) {
	id, ok := v[config.WebElementIdentifier]
	if id == "" || !ok {
		return "", fmt.Errorf(ErrorElementId, v)
	}
	return id, nil
}

func ElementsID(v []map[string]string) ([]string, error) {
	var els []string

	for _, el := range v {
		id, ok := el[config.WebElementIdentifier]
		if !ok || id == "" {
			return nil, fmt.Errorf(ErrorElementId, v)
		}
		els = append(els, id)
	}

	return els, nil
}

func (self *WebClient) addHeaders(req *http.Request, headers map[string]string) {
	if len(self.UserAgent) > 0 {
		req.Header.Set("User-Agent", self.UserAgent)
	}

	for k, v := range self.Headers {
		if _, add := headers[k]; !add {
			req.Header.Set(k, v)
		}
	}

	for _, c := range self.Cookies {
		req.AddCookie(c)
	}

	for k, v := range headers {
		if strings.ToLower(k) == "content-length" {
			if len, err := strconv.Atoi(v); err == nil && req.ContentLength <= 0 {
				req.ContentLength = int64(len)
			}
		} else if v != "" {
			req.Header.Set(k, v)
		} else {
			req.Header.Del(k)
		}
	}
}

func (self *WebClient) Request(method string, urlpath string, body io.Reader) (req *http.Request) {
	req, err := http.NewRequest(method, urlpath, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Close = self.Close
	req.Host = self.Host

	self.addHeaders(req, map[string]string{"Content-Type": "application/json"})

	return
}

func CloseResponse(r *http.Response) {
	if r != nil && r.Body != nil {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}
}

func (self *WebClient) Do(req *http.Request) (*HttpResponse, error) {
	resp, err := self.HTTPClient.Do(req)
	if urlerr, ok := err.(*url.Error); ok && urlerr.Err == errors.New("No redirect") {
		err = nil // redirect on HEAD is not an error
	}
	if err == nil {
		return &HttpResponse{*resp}, nil
	} else {
		CloseResponse(resp)
		return nil, err
	}
}

func (self *WebClient) Post(path string, content io.Reader) (*HttpResponse, error) {
	req := self.Request(http.MethodPost, path, content)
	return self.Do(req)
}

func (self *WebClient) Get(path string) (*HttpResponse, error) {
	req := self.Request(http.MethodGet, path, nil)
	return self.Do(req)
}

func (self *WebClient) Delete(path string) (*HttpResponse, error) {
	req := self.Request(http.MethodDelete, path, nil)
	return self.Do(req)
}

func (c *WebClient) Url(url, sessionId string) (*data.Url, error) {
	b := marshalData(map[string]string{"url": url})
	u := fmt.Sprintf(urlEndpoint, c.WebConfig.WebServerAddr, sessionId)
	res, err := c.Post(u, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf(ErrorOpenUrl, err)
	}

	reply := new(struct{ Value string })
	unmarshalRes(&res.Response, reply)

	return &data.Url{
		Url: reply.Value,
	}, nil
}

func (c *WebClient) Open(url, sessionId string) (*data.Url, error) {
	b := marshalData(map[string]string{"url": url})
	p := fmt.Sprintf(urlEndpoint, c.WebConfig.WebServerAddr, sessionId)

	res, err := c.Post(p, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf(ErrorOpenUrl, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value string })
	unmarshalRes(&res.Response, reply)

	return &data.Url{
		Url: reply.Value,
	}, nil
}

func (c *WebClient) NewTab(sessionId string) error {
	b := marshalData(&data.Empty{})
	p := fmt.Sprintf(newWindowEndpoint, c.WebConfig.WebServerAddr, sessionId)

	res, err := c.Post(p, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf(ErrorTab, err)
	}

	defer res.Body.Close()

	return nil
}

func (c *WebClient) Tabs(sessionId string) ([]string, error) {
	h := new(struct{ Value []string })
	url := fmt.Sprintf(windowHandlesEndpoint, c.WebConfig.WebServerAddr, sessionId)

	res, err := c.Get(url)
	if err != nil {
		return nil, fmt.Errorf(ErrorTab, err)
	}

	defer res.Body.Close()

	unmarshalRes(&res.Response, h)

	return h.Value, nil
}

func (c *WebClient) Tab(n int, sessionId string) error {
	tabs, err := c.Tabs(sessionId)
	if err != nil {
		return fmt.Errorf(ErrorTab, err)
	}

	tab := marshalData(map[string]string{"handle": tabs[n]})
	url := fmt.Sprintf(windowEndpoint, c.WebConfig.WebServerAddr, sessionId)

	res, err := c.Post(url, bytes.NewReader(tab))
	if err != nil {
		return fmt.Errorf(ErrorTab, err)
	}

	defer res.Body.Close()
	return nil
}

func (c *WebClient) FindElement(selector *data.Selector, sessionId string) (string, error) {
	body := marshalData(&data.JsonFindUsing{
		Using: selector.Using,
		Value: selector.Value,
	})

	p := fmt.Sprintf(findElementEndpoint, c.WebConfig.WebServerAddr, sessionId)
	res, err := c.Post(p, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf(ErrorFindElement, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value map[string]string })

	unmarshalRes(&res.Response, reply)
	eId, err := ElementID(reply.Value)
	if err != nil {
		if c.WebConfig.ScreenshotOnFail {
			c.Screenshot(sessionId)
		}
		return "", fmt.Errorf(ErrorElementId, reply.Value, err)
	}

	return eId, nil
}

func (c *WebClient) FindElements(selector *data.Selector, sessionId string) ([]string, error) {
	body := marshalData(&data.JsonFindUsing{
		Using: selector.Using,
		Value: selector.Value,
	})

	p := fmt.Sprintf(findElementsEndpoint, c.WebConfig.WebServerAddr, sessionId)
	res, err := c.Post(p, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf(ErrorFindElement, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value []map[string]string })

	unmarshalRes(&res.Response, reply)
	eId, err := ElementsID(reply.Value)
	if err != nil {
		if c.WebConfig.ScreenshotOnFail {
			c.Screenshot(sessionId)
		}
		return nil, fmt.Errorf(ErrorElementId, reply.Value, err)
	}

	return eId, nil
}

func (c *WebClient) FromElements(selector *data.Selector, sessionId, elementId string) ([]string, error) {
	body := marshalData(&data.JsonFindUsing{
		Using: selector.Using,
		Value: selector.Value,
	})

	p := fmt.Sprintf(fromElementsEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId)
	res, err := c.Post(p, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf(ErrorFindElement, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value []map[string]string })

	unmarshalRes(&res.Response, reply)
	eId, err := ElementsID(reply.Value)
	if err != nil {
		if c.WebConfig.ScreenshotOnFail {
			c.Screenshot(sessionId)
		}
		return nil, fmt.Errorf(ErrorElementId, reply.Value, err)
	}

	return eId, nil
}

func (c *WebClient) FromElement(selector *data.Selector, sessionId, elementId string) (string, error) {
	body := marshalData(&data.JsonFindUsing{
		Using: selector.Using,
		Value: selector.Value,
	})

	p := fmt.Sprintf(fromElementEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId)
	res, err := c.Post(p, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf(ErrorFindElement, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value map[string]string })

	unmarshalRes(&res.Response, reply)
	eId, err := ElementID(reply.Value)
	if err != nil {
		if c.WebConfig.ScreenshotOnFail {
			c.Screenshot(sessionId)
		}
		return "", fmt.Errorf(ErrorElementId, reply.Value, err)
	}

	return eId, nil
}

func (c *WebClient) TryFind(selector *data.Selector, sessionId string) (string, error) {
	body := marshalData(&data.JsonFindUsing{
		Using: selector.Using,
		Value: selector.Value,
	})

	p := fmt.Sprintf(findElementEndpoint, c.WebConfig.WebServerAddr, sessionId)
	res, err := c.Post(p, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf(ErrorFindElement, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value map[string]string })

	unmarshalRes(&res.Response, reply)
	eId, err := ElementID(reply.Value)
	if err != nil {
		return "", fmt.Errorf(ErrorFindElement, reply.Value, err)
	}

	return eId, nil
}

func (c *WebClient) Status() (*data.DriverStatus, error) {
	url := fmt.Sprintf(statusEndpoint, c.WebConfig.WebServerAddr)

	res, err := c.Get(url)
	if err != nil {
		return nil, fmt.Errorf(ErrorStatus, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value data.DriverStatus })
	unmarshalRes(&res.Response, reply)

	return &reply.Value, nil
}

func (c *WebClient) Session(caps *capabilities.Capabilities) (*data.Session, error) {
	d := marshalData(caps)

	url := fmt.Sprintf(sessionEndpoint, c.WebConfig.WebServerAddr)
	res, err := c.Post(url, bytes.NewBuffer(d))
	if err != nil {
		return nil, fmt.Errorf(ErrorCreateSession, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value data.Session })
	unmarshalRes(&res.Response, reply)

	return &data.Session{
		Id: reply.Value.Id,
	}, nil
}

func (c *WebClient) Quit(sessionId string) error {
	url := fmt.Sprintf(quitEndpoint, c.WebConfig.WebServerAddr, sessionId)
	res, err := c.Delete(url)
	if err != nil {
		return fmt.Errorf(ErrorDeleteSession, err)
	}

	defer res.Body.Close()
	return nil
}

func (c *WebClient) IsDisplayed(sessionId, elementId string) (bool, error) {
	p := fmt.Sprintf(isDisplayedEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId)
	res, err := c.Get(p)
	if err != nil {
		return false, fmt.Errorf(ErrorDisplayedElement, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value bool })
	unmarshalRes(&res.Response, reply)

	return reply.Value, nil
}

func (c *WebClient) Is(sessionId, elementId string) (bool, error) {
	p := fmt.Sprintf(isEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId)
	res, err := c.Get(p)
	if err != nil {
		if c.WebConfig.ScreenshotOnFail {
			c.Screenshot(sessionId)
		}
		return false, fmt.Errorf(ErrorDisplayedElement, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value bool })
	unmarshalRes(&res.Response, reply)

	return reply.Value, nil
}

func (c *WebClient) Click(sessionId, elementId string) error {
	p := fmt.Sprintf(clickEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId)
	d := marshalData(data.Empty{})
	res, err := c.Post(p, bytes.NewBuffer(d))
	if err != nil {
		return fmt.Errorf(ErrorClick, err)
	}

	defer res.Body.Close()

	return nil
}

func (c *WebClient) Input(keys, sessionId, elementId string) error {
	p := fmt.Sprintf(sendKeysEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId)
	d := marshalData(data.SendKeys{
		Text: keys,
	})
	res, err := c.Post(p, bytes.NewBuffer(d))
	if err != nil {
		return fmt.Errorf(ErrorSendKeys, err)
	}

	defer res.Body.Close()

	return nil
}

func (c *WebClient) Attr(attr, sessionId, elementId string) (string, error) {
	p := fmt.Sprintf(attributeEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId, attr)
	res, err := c.Get(p)
	if err != nil {
		return "", fmt.Errorf(ErrorAttribute, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value string })
	unmarshalRes(&res.Response, reply)

	return reply.Value, nil
}

func (c *WebClient) Script(script, sessionId string, args ...interface{}) error {
	if args == nil {
		args = make([]interface{}, 0)
	}

	body := marshalData(map[string]interface{}{
		"script": script,
		"args":   args,
	})

	p := fmt.Sprintf(syncScriptEndpoint, c.WebConfig.WebServerAddr, sessionId)
	res, err := c.Post(p, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf(ErrorScriptExecute, err)
	}

	defer res.Body.Close()
	return nil
}

// randSeq
// generates pseudo-random string
// for screenshot name
func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.NewSource(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (c *WebClient) Screenshot(sessionId string) error {
	data := new(struct{ Value string })

	p := fmt.Sprintf(screenshotEndpoint, c.WebConfig.WebServerAddr, sessionId)
	res, err := c.Get(p)
	if err != nil {
		return fmt.Errorf("error on screenshot request: %v", err)
	}

	unmarshalRes(&res.Response, data)

	decodedImage, err := base64.StdEncoding.DecodeString(data.Value)
	if err != nil {
		return fmt.Errorf("error on decode base64 string: %v", err)
	}

	// Create an image.Image from decoded bytes
	img, err := png.Decode(strings.NewReader(string(decodedImage)))
	if err != nil {
		return fmt.Errorf("error on decode: %v", err)
	}

	// Create a new file for the output JPEG image
	// TODO: upd randSeq, use meaninful screenshot name
	outputFile, err := os.Create(fmt.Sprintf("%s/%s_%s.jpg", config.Config.ScreenshotsPath, randSeq(8), time.Now().Format("2006_01_02_15:04:05")))
	if err != nil {
		return fmt.Errorf("error on create file: %v", err)
	}
	defer outputFile.Close()

	// Encode the image as JPEG
	err = jpeg.Encode(outputFile, img, nil)
	if err != nil {
		return fmt.Errorf("error on encode: %v", err)
	}

	return nil
}

func (c *WebClient) Active(sessionId string) (string, error) {
	p := fmt.Sprintf(activeEndpoint, c.WebConfig.WebServerAddr, sessionId)
	res, err := c.Get(p)
	if err != nil {
		return "", fmt.Errorf(ErrorActiveElement, err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value map[string]string })

	unmarshalRes(&res.Response, reply)
	eId, err := ElementID(reply.Value)
	if err != nil {
		if c.WebConfig.ScreenshotOnFail {
			c.Screenshot(sessionId)
		}
		return "", fmt.Errorf(ErrorActiveElement, reply.Value, err)
	}

	return eId, nil
}

func (c *WebClient) Action(keys, action, sessionId string) error {
	actions := make([]data.KeyAction, 0, len(keys))

	for _, key := range keys {
		actions = append(actions, data.KeyAction{
			Type: action,
			Key:  string(key),
		})
	}
	p := fmt.Sprintf(actionEndpoint, c.WebConfig.WebServerAddr, sessionId)

	data := marshalData(map[string]interface{}{
		"actions": []interface{}{
			map[string]interface{}{
				"type":    "key",
				"id":      "default keyboard",
				"actions": actions,
			}},
	})

	res, err := c.Post(p, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf(ErrorAction, err)
	}

	defer res.Body.Close()

	return nil
}

func (c *WebClient) ReleaseAction(sessionId string) error {
	p := fmt.Sprintf(actionEndpoint, c.WebConfig.WebServerAddr, sessionId)

	res, err := c.Delete(p)
	if err != nil {
		return fmt.Errorf(ErrorActiveElement, err)
	}

	defer res.Body.Close()

	return nil
}

func (c *WebClient) Text(sessionId, elementId string) (string, error) {
	p := fmt.Sprintf(textEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId)

	res, err := c.Get(p)
	if err != nil {
		return "", fmt.Errorf(ErrorTextElement, err)
	}

	defer res.Body.Close()

	t := new(struct{ Value string })
	unmarshalRes(&res.Response, t)

	return t.Value, nil
}
