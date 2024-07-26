package test

import (
	"testing"

	"github.com/mcsymiv/gost/gost"
)

func TestStep(t *testing.T) {
	st := gost.New(t)

	defer st.Tear()

	st.Open("http://google.com")
	st.Type("hello", "//*[@id='APjFqb']")
}
