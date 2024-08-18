package test

import (
	"os"

	"github.com/mcsymiv/gost/driver"
	"github.com/xlzd/gotp"
)

func loginOkta(d *driver.WebDriver) {
	d.F("//*[@id='okta-signin-username']").Input(os.Getenv("OKTA_LOGIN"))
	d.F("//*[@id='okta-signin-password']").Input(os.Getenv("OKTA_PASS")).Input(driver.EnterKey)
	totp := gotp.NewDefaultTOTP(os.Getenv("OKTA_TOTP"))
	d.F("//*[@id='input59']").Input(totp.Now()).Input(driver.EnterKey)
}
