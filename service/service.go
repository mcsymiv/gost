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

	sm.HandleFunc("GET /hello", wd.get())
	sm.HandleFunc("GET /status", wd.get())
	sm.HandleFunc("POST /session", wd.post())
	sm.HandleFunc("DELETE /session/{sessionId}", wd.delete())
	sm.HandleFunc("POST /session/{sessionId}/url", wd.post())
	sm.HandleFunc("POST /session/{sessionId}/element", wd.post())

	return sm
}

func (wd *WebDriverHandler) post() func(http.ResponseWriter, *http.Request) {
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

func (wd *WebDriverHandler) get() func(http.ResponseWriter, *http.Request) {
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

func (wd *WebDriverHandler) delete() func(http.ResponseWriter, *http.Request) {
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
