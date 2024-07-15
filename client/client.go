package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/data"
)

const (
	sessionEndpoint     = "%s/session"
	quitEndpoint        = "%s/session/%s"
	urlEndpoint         = "%s/session/%s/url"
	findElementEndpoint = "%s/session/%s/element"
	isDisplayedEndpoint = "%s/session/%s/element/%s/displayed"
	clickEndpoint       = "%s/session/%s/element/%s/click"
)

const (
	// LegacyWebElementIdentifier is the string constant used in the old Selenium 2 protocol
	// WebDriver JSON protocol that is the key for the map that contains an
	// unique element identifier.
	// This value is ignored in element id retreival
	LegacyWebElementIdentifier = "ELEMENT"

	// WebElementIdentifier is the string constant defined by the W3C Selenium 3 protocol
	// specification that is the key for the map that contains a unique element identifier.
	WebElementIdentifier = "element-6066-11e4-a52e-4f735466cecf"

	// ShadowRootIdentifier A shadow root is an abstraction used to identify a shadow root when
	// it is transported via the protocol, between remote and local ends.
	ShadowRootIdentifier = "shadow-6066-11e4-a52e-4f735466cecf"
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
				delay:      time.Duration(200 * time.Millisecond),
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
				delay:      time.Duration(200 * time.Millisecond),
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
	id, ok := v[WebElementIdentifier]
	if !ok || id == "" {
		panic(fmt.Sprintf("Error on find element: %v", v))
	}
	return id
}

func ElementsID(v []map[string]string) []string {
	var els []string

	for _, el := range v {
		id, ok := el[WebElementIdentifier]
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

func (self *WebClient) NewRequest(method string, urlpath string, body io.Reader, headers map[string]string) (req *http.Request) {
	req, err := http.NewRequest(strings.ToUpper(method), urlpath, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Close = self.Close
	req.Host = self.Host

	self.addHeaders(req, headers)

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

// ParamValues fills the input url.Values according to params
func ParamValues(params map[string]interface{}, q url.Values) url.Values {
	if q == nil {
		q = url.Values{}
	}

	for k, v := range params {
		val := reflect.ValueOf(v)

		switch val.Kind() {
		case reflect.Slice:
			if val.IsNil() { // TODO: add an option to ignore empty values
				q.Set(k, "")
				continue
			}
			fallthrough

		case reflect.Array:
			for i := 0; i < val.Len(); i++ {
				av := val.Index(i)
				q.Add(k, fmt.Sprintf("%v", av))
			}

		default:
			q.Set(k, fmt.Sprintf("%v", v))
		}
	}

	return q
}

func URLWithParams(base string, params map[string]interface{}) (u *url.URL) {
	return URLWithPathParams(base, "", params)
}

// Given a base URL and a bag of parameteters returns the URL with the encoded parameters
func URLWithPathParams(base string, path string, params map[string]interface{}) (u *url.URL) {
	u, err := url.Parse(base)
	if err != nil {
		log.Fatal(err)
	}

	if len(path) > 0 {
		u, err = u.Parse(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	q := ParamValues(params, u.Query())
	u.RawQuery = q.Encode()
	return u
}

func (self *WebClient) Post(path string, content io.Reader) (*HttpResponse, error) {
	req := self.Request(http.MethodPost, path, content)
	return self.Do(req)
}

func (self *WebClient) GetWithParams(path string, params map[string]interface{}) (*HttpResponse, error) {
	req := self.Request(http.MethodGet, URLWithParams(path, params).String(), nil)
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
