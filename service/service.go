package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func Recoverer(next http.Handler) http.Handler {
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

func Retry(next http.Handler) http.Handler {
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
	conf *config.WebConfig
}

func Handler() http.Handler {
	sm := http.NewServeMux()

	wd := &WebDriverHandler{
		conf: config.Config,
	}

	sm.HandleFunc("GET /hello", wd.Hello())
	sm.HandleFunc("GET /status", wd.DriverStatus())
	sm.HandleFunc("POST /session", wd.CreateSession())
	sm.HandleFunc("DELETE /session", wd.Quit())
	sm.HandleFunc("POST /session/{sessionId}/url", wd.Url())
	sm.HandleFunc("POST /session/{sessionId}/element", wd.FindElement())

	return sm
}

func (wd *WebDriverHandler) Url() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		fmt.Println(url)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read session request body: %v", err))
			return
		}

		res, err := http.Post(url, config.ApplicationJson, bytes.NewBuffer(data))
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on post session request: %v", err))
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read session response: %v", err))
			return
		}

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(body)
	}
}

func (wd *WebDriverHandler) Quit() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, req.URL.Path)
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

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(data)
	}
}

func (wd *WebDriverHandler) CreateSession() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		fmt.Println(url)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read session request body: %v", err))
			return
		}

		res, err := http.Post(url, config.ApplicationJson, bytes.NewBuffer(data))
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on post session request: %v", err))
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read session response: %v", err))
			return
		}

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(body)
	}
}

func (wd *WebDriverHandler) Hello() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello gost"))
	}
}

func (wd *WebDriverHandler) DriverStatus() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		res, err := http.Get(url)
		if err != nil {
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read session response: %v", err))
			return
		}

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(data)
	}
}

func (wd *WebDriverHandler) FindElement() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read url request body: %v", err))
			return
		}

		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		res, err := http.Post(url, config.ApplicationJson, bytes.NewBuffer(body))
		if err != nil {
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return
		}

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(data)
	}
}
