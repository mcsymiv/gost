package test

import (
	"testing"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/gost"
)

func TestDriver(t *testing.T) {
	st := gost.Gost(
		t,
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer st.Shutdown()

	st.Open("https://google.com").
		Click("//*[@id='APjFqb']")
}
