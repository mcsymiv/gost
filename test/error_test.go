package test

import (
	"testing"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/gost"
)

func TestError(t *testing.T) {
	d, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)

	defer tear()

	d.Open("http://192.168.0.1/")
	d.F("//*[@type='password']").Is()

}
