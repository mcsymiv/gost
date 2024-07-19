package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mcsymiv/gost/config"
)

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

func (wd *WebDriverHandler) postRetrier(v verifier, next http.Handler) http.Handler {

	var res *http.Response

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		data, err := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(data))
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Errorf("error on read post request body: %v", err))
		}
		start := time.Now()
		end := start.Add(30 * time.Second)

		for {
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
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

func (wd *WebDriverHandler) getRetrier(v verifier, next http.Handler) http.Handler {

	var res *http.Response

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, r.URL.Path)
		start := time.Now()
		end := start.Add(30 * time.Second)

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
			if v.verify(res) {
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

func (wd *WebDriverHandler) isRetrier(v verifier, next http.Handler) http.Handler {

	var isDisplayed struct{ Value bool }
	var res *http.Response

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		servicePath := r.URL.Path
		wdPath := strings.Replace(servicePath, "is", "displayed", 1)

		fmt.Println(wdPath)
		url := fmt.Sprintf("%s%s", wd.conf.WebDriverAddr, wdPath)
		start := time.Now()
		end := start.Add(10 * time.Second)

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
			isDisplayed.Value = v.verify(res)
			if isDisplayed.Value {
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

		body, err := json.Marshal(isDisplayed)
		if err != nil {
			fmt.Println("errr")
			json.NewEncoder(w).Encode(fmt.Errorf("error on read post response: %v", err))
			return
		}

		defer res.Body.Close()

		w.Header().Set(config.ContenType, config.ApplicationJson)
		w.Write(body)

		next.ServeHTTP(w, r)
	})
}
