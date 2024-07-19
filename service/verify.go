package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type verifyStatusOk struct {
}

type verifyDisplay struct {
}

// type requestVerifier interface {
// 	verify(*http.Response, ...interface{}) bool
// }

type verifier interface {
	verify(*http.Response) bool
}

func (v *verifyStatusOk) verify(res *http.Response) bool {
	if res.StatusCode == http.StatusOK {
		return true
	}

	return false
}

// verify isDisplayStreategy
// will assign true to b to reuse in IsDisplayed()
func (v *verifyDisplay) verify(res *http.Response) bool {
	if res.StatusCode == http.StatusOK {
		var displayResponse = new(struct{ Value bool })

		err := json.NewDecoder(res.Body).Decode(displayResponse)
		if err != nil {
			fmt.Println("error on json NewDecoder")
			res.Body.Close()
			panic(err)
		}

		if displayResponse.Value {
			res.Body.Close()
			return true
		}

		res.Body.Close()
		return false
	}

	return false
}
