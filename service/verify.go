package service

import (
	"encoding/json"
	"log"
	"net/http"
)

type verifyStatusOk struct {
	Method string
}

type verifyDisplay struct {
	Method string
}

type requestVerifier interface {
	verify(*http.Response, ...interface{}) bool
	method() string
}

func (v *verifyStatusOk) method() string {
	return v.Method
}

func (v *verifyDisplay) method() string {
	return v.Method
}

func (v *verifyStatusOk) verify(res *http.Response, b ...interface{}) bool {
	if res.StatusCode == http.StatusOK {
		return true
	}

	return false
}

// verify isDisplayStreategy
// will assign true to b to reuse in IsDisplayed()
func (v *verifyDisplay) verify(res *http.Response, b ...interface{}) bool {
	if res.StatusCode == http.StatusOK {
		var displayResponse = new(struct{ Value bool })

		err := json.NewDecoder(res.Body).Decode(displayResponse)
		if err != nil {
			log.Println("error on json NewDecoder")
			res.Body.Close()
			panic(err)
		}

		if displayResponse.Value {
			b[0] = true
			res.Body.Close()
			return true
		}

		res.Body.Close()
		return false
	}

	return false
}
