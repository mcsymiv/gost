package test

import (
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/gost"
)

func TestScript(t *testing.T) {
	d, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)

	defer tear()

	d.Open("http://192.168.0.1/")
	d.F("//*[@type='password']").Keys(os.Getenv("HOME_PASS"))

	id := d.F("LOG IN").Id()

	d.ExecuteScript("click", id)
	time.Sleep(7 * time.Second)
	d.F("Clients").Click()

	time.Sleep(7 * time.Second)
}
