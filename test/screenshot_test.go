package test

import (
	"testing"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/gost"
)

func TestScreenshot(t *testing.T) {
	d, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)

	defer tear()

	d.Open("https://google.com")
	d.Screenshot()
}
