package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/gost"
	"github.com/xlzd/gotp"
)

func loginOkta(d *driver.WebDriver) {
	d.F("//*[@id='okta-signin-username']").Input(os.Getenv("OKTA_LOGIN"))
	d.F("//*[@id='okta-signin-password']").Input(os.Getenv("OKTA_PASS")).Input(driver.EnterKey)
	totp := gotp.NewDefaultTOTP(os.Getenv("OKTA_TOTP"))
	d.F("//*[@id='input59']").Input(totp.Now()).Input(driver.EnterKey)
}

// set location attr on new account
func attribute() {
	d, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	d.Open(os.Getenv("BG_ENV_QA_DEV_01"))
	d.Cl("Elateral SSO")
	loginOkta(d)
	d.Cl("My Account")
	d.Cl("Manage Attributes")
	d.Cl("10")
	d.Cl("200")
	d.Cl("//*[text()='Location']/..//*[@data-qa-id='gears']")
	fmt.Println(d.F("Name *").Attr("for"))
	d.Cl("Add Root Item")
	d.F("Add Item").P(2).Next("Name *").Click()
	d.Active().Input("ae")
	d.Cl("Add")
	time.Sleep(time.Second * 5)
}

func TestBg(t *testing.T) {
	args := os.Args[6:]
	switch args[0] {
	case "new":
		attribute()
	}
}
