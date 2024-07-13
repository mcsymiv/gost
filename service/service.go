package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mcsymiv/gost/config"
)

type WebServer struct {
	Server        *http.Server
	WebDriverAddr string
}

func NewServer(h http.Handler) *WebServer {
	return &WebServer{
		Server: &http.Server{
			Addr:    ":8080",
			Handler: h,
		},
	}
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

type WebDriverHandler struct {
	conf   *config.WebConfig
	client *http.Client
}

func Handler() http.Handler {
	sm := http.NewServeMux()

	wd := &WebDriverHandler{
		conf:   config.Config,
		client: &http.Client{},
	}

	sm.HandleFunc("GET /hello", wd.get())
	sm.Handle("GET /status", logger(http.HandlerFunc(wd.get())))
	sm.HandleFunc("POST /session", wd.post())
	sm.HandleFunc("DELETE /session/{sessionId}", wd.delete())
	sm.HandleFunc("POST /session/{sessionId}/url", wd.post())
	sm.Handle("POST /session/{sessionId}/element", wd.retry(&verifyStatusOk{http.MethodPost}, wd.post()))
	// TODO: add display_test
	sm.Handle("GET /session/{sessionId}/element/{elementId}/displayed", wd.retry(&verifyDisplay{http.MethodGet}, wd.get(), new(struct{ Value bool })))

	return sm
}

func (wd *WebDriverHandler) retry(v requestVerifier, next http.Handler, b ...interface{}) http.Handler {

	var res *http.Response

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read post request body: %v", err))
		}
		start := time.Now()
		end := start.Add(30 * time.Second)

		for {
			req, err := http.NewRequest(v.method(), url, bytes.NewBuffer(data))
			if err != nil {
				fmt.Println("error on NewRequest")
				req.Body.Close()
				panic(err)
			}

			res, err = wd.client.Do(req)
			if err != nil {
				fmt.Println("error on Client Do Request")
				res.Body.Close()
				panic(err)
			}

			// strategy for strategy
			// "verified" response will return true
			// and break out of the loop
			if v.verify(res, b) {
				break
			}

			// close res res.Body if not verified
			// i.e. loopStrategyRequest returns false
			res.Body.Close()

			if time.Now().After(end) {
				log.Println("timeout")
				// if config.TestSetting.ScreenshotOnFail {
				if wd.conf.ScreenshotOnFail {
					fmt.Println("screnshot")
					// d.Screenshot()
				}

				break
			}

			time.Sleep(300 * time.Millisecond)
			fmt.Println("retry find element")
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read post response: %v", err))
			return
		}

		defer res.Body.Close()

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(body)

		next.ServeHTTP(w, r)
	})
}

func (wd *WebDriverHandler) post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read post request body: %v", err))
			return
		}

		res, err := http.Post(url, config.ApplicationJson, bytes.NewBuffer(data))
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on post request: %v", err))
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read post response: %v", err))
			return
		}
		defer res.Body.Close()

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(body)
	}
}

func (wd *WebDriverHandler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		res, err := http.Get(url)
		if err != nil {
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on get response: %v", err))
			return
		}
		defer res.Body.Close()

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(data)
	}
}

// func (wd *WebDriverHandler) get3(w http.ResponseWriter, r *http.Request) ([]byte, error) {
// 	url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
// 	res, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	data, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		json.NewEncoder(w).Encode(fmt.Errorf("error on get response: %v", err))
// 		return nil, err
// 	}
// 	defer res.Body.Close()
//
// 	return data, nil
// }
//
// func (wd *WebDriverHandler) get2() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		data, err := wd.get3(w, r)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
//
// 		w.Header().Set(config.ContenType, config.ApplicationJson)
// 		w.Write(data)
// 	}
// }

func (wd *WebDriverHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		wdReq, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			return
		}

		res, err := http.DefaultClient.Do(wdReq)
		if err != nil {
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read session response: %v", err))
			return
		}
		defer res.Body.Close()

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(data)
	}
}
