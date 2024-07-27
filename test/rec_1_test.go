package test

import (
	"testing"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/gost"
)

func TestRecord1(t *testing.T) {
	st := gost.New(t, capabilities.BrowserName("chrome"))
	defer st.Tear()

	st.Open("chrome://new-tab-page/")

	st.Open("https://www.google.com/")

	st.Click("//*[@id='APjFqb']")

	st.Keys("hello")

}
