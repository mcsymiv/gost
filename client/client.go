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

const (
	// W3C Endpoints
	sessionEndpoint     = "%s/session"
	quitEndpoint        = "%s/session/%s"
	urlEndpoint         = "%s/session/%s/url"
	findElementEndpoint = "%s/session/%s/element"
	isDisplayedEndpoint = "%s/session/%s/element/%s/displayed"
	clickEndpoint       = "%s/session/%s/element/%s/click"
	sendKeysEndpoint    = "%s/session/%s/element/%s/value"
	attributeEndpoint   = "%s/session/%s/element/%s/attribute/%s"
	screenshotEndpoint  = "%s/session/%s/screenshot"

	// GoST Endpoints
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

func WebDriverClient() *WebClient {
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

func ElementID(v map[string]string) string {
	id, ok := v[config.WebElementIdentifier]
	if !ok || id == "" {
		panic(fmt.Sprintf("Error on find element: %v", v))
	}
	return id
}

func ElementsID(v []map[string]string) []string {
	var els []string

	for _, el := range v {
		id, ok := el[config.WebElementIdentifier]
		if !ok || id == "" {
			panic(fmt.Sprintf("Error on find elements: %v", v))
		}
		els = append(els, id)
	}

	return els
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
		return nil, fmt.Errorf("error on url request: %v", err)
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
		return nil, fmt.Errorf("error on open request: %v", err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value string })
	unmarshalRes(&res.Response, reply)

	return &data.Url{
		Url: reply.Value,
	}, nil
}

func (c *WebClient) FindElement(selector *data.Selector, sessionId string) (string, error) {
	body := marshalData(&data.JsonFindUsing{
		Using: selector.Using,
		Value: selector.Value,
	})

	p := fmt.Sprintf(findElementEndpoint, c.WebConfig.WebServerAddr, sessionId)
	res, err := c.Post(p, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error on find element request: %v", err)
	}

	defer res.Body.Close()

	reply := new(struct{ Value map[string]string })
	unmarshalRes(&res.Response, reply)
	eId := ElementID(reply.Value)

	return eId, nil
}

func (c *WebClient) Status() (*data.DriverStatus, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/status", c.WebConfig.WebServerAddr), nil)
	if err != nil {
		return nil, fmt.Errorf("error on new status request: %v", err)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on status request: %v", err)
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
		return nil, fmt.Errorf("error on session request: %v", err)
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
		return fmt.Errorf("error on delete session request: %v", err)
	}

	defer res.Body.Close()
	return nil
}

func (c *WebClient) IsDisplayed(sessionId, elementId string) (bool, error) {
	p := fmt.Sprintf(isDisplayedEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId)
	res, err := c.Get(p)
	if err != nil {
		return false, fmt.Errorf("error on find element request: %v", err)
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
		fmt.Println("errr")
		return false, fmt.Errorf("error on is request: %v", err)
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
		return fmt.Errorf("error on click request: %v", err)
	}

	defer res.Body.Close()

	return nil
}

func (c *WebClient) Keys(keys, sessionId, elementId string) error {
	p := fmt.Sprintf(sendKeysEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId)
	d := marshalData(data.SendKeys{
		Text: keys,
	})
	res, err := c.Post(p, bytes.NewBuffer(d))
	if err != nil {
		return fmt.Errorf("error on keys request: %v", err)
	}

	defer res.Body.Close()

	return nil
}

func (c *WebClient) Attr(attr, sessionId, elementId string) (string, error) {
	p := fmt.Sprintf(attributeEndpoint, c.WebConfig.WebServerAddr, sessionId, elementId, attr)
	res, err := c.Get(p)
	if err != nil {
		return "", fmt.Errorf("error on attribute request: %v", err)
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
		return fmt.Errorf("error on find element request: %v", err)
	}

	defer res.Body.Close()
	return nil
}

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

	fmt.Println("pathh", config.Config.ScreenshotsPath)
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
