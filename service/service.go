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
	sm.Handle("GET /status", logger(wd.get()))
	sm.Handle("POST /session", logger(wd.post()))
	sm.HandleFunc("DELETE /session/{sessionId}", wd.delete())
	sm.HandleFunc("POST /session/{sessionId}/url", wd.post())

	sm.Handle("POST /session/{sessionId}/element", logger(wd.retrier(&verifyStatusOk{})))
	sm.Handle("POST /session/{sessionId}/elements", logger(wd.retrier(&verifyStatusOk{})))
	sm.Handle("POST /session/{sessionId}/element/{elementId}/element", logger(wd.retrier(&verifyStatusOk{})))
	sm.Handle("POST /session/{sessionId}/element/{elementId}/elements", logger(wd.retrier(&verifyStatusOk{})))
	sm.HandleFunc("GET /session/{sessionId}/element/active", wd.get())

	sm.Handle("POST /session/{sessionId}/element/{elementId}/click", wd.retrier(&verifyStatusOk{}))
	sm.HandleFunc("POST /session/{sessionId}/element/{elementId}/value", wd.post())
	sm.HandleFunc("GET /session/{sessionId}/element/{elementId}/text", wd.get())
	sm.Handle("GET /session/{sessionId}/element/{elementId}/displayed", wd.isRetrier(&verifyValue{}))
	sm.Handle("GET /session/{sessionId}/element/{elementId}/is", wd.isDisplayed(wd.isRetrier(&verifyValue{})))
	sm.Handle("GET /session/{sessionId}/element/{elementId}/attribute/{attribute}", wd.retrier(&verifyStatusOk{}))

	sm.Handle("POST /session/{sessionId}/script", wd.script(wd.post()))
	sm.HandleFunc("GET /session/{sessionId}/screenshot", wd.get())

	sm.HandleFunc("POST /session/{sessionId}/window", wd.post())
	sm.HandleFunc("GET /session/{sessionId}/window/handles", wd.get())
	sm.HandleFunc("POST /session/{sessionId}/window/new", wd.post())

	sm.HandleFunc("POST /session/{sessionId}/actions", wd.post())
	sm.HandleFunc("DELETE /session/{sessionId}/actions", wd.delete())
	return sm
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
