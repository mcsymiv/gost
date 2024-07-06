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
	"github.com/mcsymiv/gost/models"
)

// RestClient represents a REST client configuration.
type WebClient struct {
	Close          bool
	Host           string
	WebServerAddr2 *url.URL
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
	if self.WebServerAddr2 != nil {
		if u, err := self.WebServerAddr2.Parse(urlpath); err != nil {
			log.Fatal(err)
		} else {
			urlpath = u.String()
		}
	}

	req, err := http.NewRequest(strings.ToUpper(method), urlpath, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Close = self.Close
	req.Host = self.Host

	self.addHeaders(req, map[string]string{"Content-Type": "application/json"})

	return
}

func (self *WebClient) NewRequest(method string, urlpath string, body io.Reader, headers map[string]string) (req *http.Request) {
	if self.WebServerAddr2 != nil {
		if u, err := self.WebServerAddr2.Parse(urlpath); err != nil {
			log.Fatal(err)
		} else {
			urlpath = u.String()
		}
	}

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
	req := self.Request("POST", path, content)
	return self.Do(req)
}

func (self *WebClient) GetWithParams(path string, params map[string]interface{}) (*HttpResponse, error) {
	req := self.Request("GET", URLWithParams(path, params).String(), nil)
	return self.Do(req)
}

func (self *WebClient) Get(path string) (*HttpResponse, error) {
	req := self.Request("GET", path, nil)
	return self.Do(req)
}

func (self *WebClient) Delete(path string) (*HttpResponse, error) {
	req := self.Request("DELETE", path, nil)
	return self.Do(req)
}

func (c *WebClient) Url(url string) (*models.Url, error) {
	data := marshalData(map[string]string{"url": url})
	u := fmt.Sprintf("%s/url", c.WebConfig.WebServerAddr)
	res, err := c.Post(u, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error on url request: %v", err)
	}

	reply := new(struct{ Value string })
	unmarshalRes(&res.Response, reply)

	return &models.Url{
		Url: reply.Value,
	}, nil
}

func (c *WebClient) Open(url, sessionId string) (*models.Url, error) {
	data := marshalData(map[string]string{"url": url})
	p := fmt.Sprintf("%s/session/%s/url", c.WebConfig.WebServerAddr, sessionId)

	res, err := c.Post(p, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error on open request: %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error on read open body: %v", err)
	}

	reply := new(struct{ Value string })
	if err := json.Unmarshal(body, reply); err != nil {
		return nil, fmt.Errorf("erron on open unmarshal: %v", err)
	}

	return &models.Url{
		Url: reply.Value,
	}, nil
}

func (c *WebClient) FindElement(sessionId string) (*models.Url, error) {
	data := marshalData(map[string]string{"url": ""})
	p := fmt.Sprintf("%s/session/%s/url", c.WebConfig.WebServerAddr, sessionId)

	res, err := c.Post(p, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error on open request: %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error on read open body: %v", err)
	}

	reply := new(struct{ Value string })
	if err := json.Unmarshal(body, reply); err != nil {
		return nil, fmt.Errorf("erron on open unmarshal: %v", err)
	}

	return &models.Url{
		Url: reply.Value,
	}, nil
}

func (c *WebClient) Status() (*models.DriverStatus, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/status", c.WebConfig.WebServerAddr), nil)
	if err != nil {
		return nil, fmt.Errorf("error on new status request: %v", err)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on status request: %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error on read status body: %v", err)
	}

	reply := new(struct{ Value models.DriverStatus })
	if err := json.Unmarshal(body, reply); err != nil {
		return nil, fmt.Errorf("error on unmarshal status body: %v", err)
	}

	return &reply.Value, nil
}

func (c *WebClient) Session(caps *capabilities.Capabilities) (*models.Session, error) {
	data, err := json.Marshal(caps)
	if err != nil {
		return nil, fmt.Errorf("error on new driver marshal caps: %v", err)
	}

	url := fmt.Sprintf("%s/session", c.WebConfig.WebServerAddr)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error on new request: %v", err)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on session request: %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error on read session body: %v", err)
	}

	reply := new(struct{ Value models.Session })
	if err := json.Unmarshal(body, reply); err != nil {
		return nil, fmt.Errorf("erron on session unmarshal: %v", err)
	}

	return &models.Session{
		Id: reply.Value.Id,
	}, nil
}
