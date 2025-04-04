package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/mcsymiv/gost/config"
)

func sessionId(url string) string {
	var sessionUUID = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	return sessionUUID.FindString(url)
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

func (wd *WebDriverHandler) retrier(v verifier) http.Handler {
	var res *http.Response

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)

		var data []byte
		var err error

		if r.Body != http.NoBody {
			data, err = io.ReadAll(r.Body)
			if err != nil {
				json.NewEncoder(w).Encode(fmt.Errorf("error on read post request body: %v", err))
			}

			r.Body = io.NopCloser(bytes.NewReader(data))
		}

		start := time.Now()
		end := start.Add(wd.conf.WaitForTimeout * time.Second)

		for {
			req, err := http.NewRequest(r.Method, url, bytes.NewReader(data))
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
			if v.verify(res) {
				break
			}

			if time.Now().After(end) {
				break
			}

			time.Sleep(wd.conf.WaitForInterval * time.Millisecond)
			fmt.Println("retry find element")

			// close res res.Body if not verified
			// i.e. StrategyRequest returns false
			res.Body.Close()
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on get response: %v", err))
			return
		}

		defer res.Body.Close()

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(data)
	})
}

func (wd *WebDriverHandler) isRetrier(v verifier) http.Handler {
	var ok struct{ Value bool }
	var res *http.Response

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		start := time.Now()
		end := start.Add(wd.conf.WaitForTimeout * time.Second)

		for {
			req, err := http.NewRequest(http.MethodGet, url, nil)
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
			ok.Value = v.verify(res)
			if ok.Value {
				break
			}

			if time.Now().After(end) {
				break
			}

			time.Sleep(wd.conf.WaitForInterval * time.Millisecond)
			fmt.Println("retry on displayed")

			// close res res.Body if not verified
			// i.e. StrategyRequest returns false
			res.Body.Close()
		}

		body, err := json.Marshal(ok)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read post response: %v", err))
			return
		}

		defer res.Body.Close()

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(body)
	})
}

func (wd *WebDriverHandler) isDisplayed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		servicePath := r.URL.Path
		wdPath := strings.Replace(servicePath, "is", "displayed", 1)

		r.URL.Path = wdPath

		next.ServeHTTP(w, r)
	})
}

func (wd *WebDriverHandler) script(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		servicePath := r.URL.Path
		wdPath := strings.Replace(servicePath, "script", "execute/sync", 1)

		r.URL.Path = wdPath

		next.ServeHTTP(w, r)
	})
}

func (wd *WebDriverHandler) retrier2(v verifier) http.Handler {
	var res *http.Response

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)

		var data []byte
		var err error

		if r.Body != http.NoBody {
			data, err = io.ReadAll(r.Body)
			if err != nil {
				json.NewEncoder(w).Encode(fmt.Errorf("error on read post request body: %v", err))
			}

			r.Body = io.NopCloser(bytes.NewReader(data))
		}

		start := time.Now()
		end := start.Add(wd.conf.WaitForTimeout * time.Second)

		for {
			req, err := http.NewRequest(r.Method, url, bytes.NewReader(data))
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
			if v.verify(res) {
				break
			}

			if time.Now().After(end) {
				break
			}

			time.Sleep(wd.conf.WaitForInterval * time.Millisecond)
			fmt.Println("retry find element")

			// close res res.Body if not verified
			// i.e. StrategyRequest returns false
			res.Body.Close()
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on get response: %v", err))
			return
		}

		defer res.Body.Close()

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(data)
	})
}
