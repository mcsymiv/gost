package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/gost"
)

func token() {
	d, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	d.Url("https://console.cloud.google.com/")
	d.F("//*[@id='identifierId']").Input(os.Getenv("G_USER")).Input(driver.EnterKey)
	time.Sleep(2 * time.Second)
	d.F("Enter your password").Input(os.Getenv("G_PASS")).Input(driver.EnterKey)
	time.Sleep(1 * time.Second)
	d.ClickJs("//*[text()='Next']")
	d.Cl("APIs and services")
	d.Cl(" Credentials ")
	d.Cl("//services-oauth-clients-table-oauth-clients-gql//*[@data-mat-icon-name='delete']")
	d.Cl("Delete")
	d.ClickJs("//*[contains(text()'Create credentials')]")
	d.Cl(" OAuth client ID ")
	d.Cl("Application type")
	d.Cl("Desktop app")
	d.Cl("//button//span[contains(text(),'Create')]")
	d.Cl("Download JSON")

	time.Sleep(4 * time.Second)
}

func newtoken() {
	d, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	d.Url("https://console.cloud.google.com/")
	d.F("//*[@id='identifierId']").Input(os.Getenv("G_USER")).Input(driver.EnterKey)
	time.Sleep(2 * time.Second)
	d.F("Enter your password").Input(os.Getenv("G_PASS")).Input(driver.EnterKey)
	time.Sleep(1 * time.Second)
	d.ClickJs("//*[text()='Next']")
	d.Cl("APIs and services")
	d.Cl(" Credentials ")
	time.Sleep(1 * time.Second)
	d.Cl("//*[contains(text(),'Create credentials')]")
	d.Cl(" OAuth client ID ")
	d.Cl("Application type")
	d.Cl("Desktop app")
	d.Cl("//button//span[contains(text(),'Create')]")
	d.Cl("Download JSON")

	time.Sleep(4 * time.Second)
}

func TestGc(t *testing.T) {
	args := os.Args[6:]

	if len(args) != 0 {
		switch args[0] {
		case "token":
			token()
		case "new":
			newtoken()
		default:
			fmt.Println(`
				Usage:
				make gc token
				make gc new
			`)
		}
	}
}
