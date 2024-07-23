package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/gost"
)

func TestScript(t *testing.T) {
	d, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)

	defer tear()

	d.Open("http://192.168.0.1/")
	d.F("//*[@type='password']").Keys(os.Getenv("HOME_PASS"))

	f, err := config.FindFile("../js", "click.js")
	if err != nil {
		panic(fmt.Sprintf("error on find file: %v", err))
	}
	c, err := os.ReadFile(f)
	if err != nil {
		panic(fmt.Sprintf("error on read file: %v", err))
	}

	id := d.F("LOG IN").ElementIdentifier()

	d.Script(string(c), id)
	time.Sleep(7 * time.Second)
	d.F("Clients").Click()

	time.Sleep(7 * time.Second)
}
