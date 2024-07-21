package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type verifyStatusOk struct{}

type verifyValue struct{}

// type requestVerifier interface {
// 	verify(*http.Response, ...interface{}) bool
// }

type verifier interface {
	verify(*http.Response) bool
}

// verify
func (v *verifyStatusOk) verify(res *http.Response) bool {
	if res.StatusCode == http.StatusOK {
		return true
	}

	return false
}

// verify
func (v *verifyValue) verify(res *http.Response) bool {
	if res.StatusCode == http.StatusOK {
		var okValue = new(struct{ Value bool })

		err := json.NewDecoder(res.Body).Decode(okValue)
		if err != nil {
			fmt.Println("error on json NewDecoder")
			res.Body.Close()
			panic(err)
		}

		if okValue.Value {
			res.Body.Close()
			return true
		}

		res.Body.Close()
		return false
	}

	return false
}
