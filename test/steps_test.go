package test

import (
	"testing"

	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/gost"
)

func TestStep(t *testing.T) {
	st := gost.New(t)

	defer st.Tear()

	st.Open("http://google.com")
	st.Type("hello", "//*[@id='APjFqb']")
}

func TestUntil(t *testing.T) {
	st := gost.New(t)

	defer st.Tear()

	st.Open("http://google.com")
	st.Keys("hello")
	st.Keys(driver.EnterKey)

	st.Until(func() bool {
		return st.Is("//*[@id='APjFqb']")
	})
}
