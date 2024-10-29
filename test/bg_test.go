package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/gost"
	"github.com/xlzd/gotp"
)

var d *driver.WebDriver

func TestBg(t *testing.T) {
	dr, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)

	d = dr
	defer tear()

	args := os.Args[6:]
	runTest(args[0])
}

func runTest(testName string) {
	var steps map[string][]func() = map[string][]func(){
		"1": []func(){
			open(os.Getenv("01"), os.Getenv("DOMAIN")),
			login(os.Getenv("PROF")),
		},

		"2": []func(){
			open(os.Getenv("01"), os.Getenv("DOMAIN")),
			login(os.Getenv("PROF")),
			option("Manage UI"),
		},

		"3": []func(){
			open(os.Getenv("01"), os.Getenv("DOMAIN")),
			login(os.Getenv("PROF")),
			option("Manage Attributes"),
			addAttribute("Attr Name", "List"),
		},
	}

	// before all
	for _, fn := range steps[testName] {
		// before each
		fn()
		// after each
	}
	// after all
}

func open(env, name string) func() {
	return func() {
		d.Open(fmt.Sprintf("https://%s.%s", env, name))
	}
}

func login(prof string) func() {
	return func() {
		d.Cl(prof)
		d.F("//*[@id='okta-signin-username']").Input(os.Getenv("OKTA_LOGIN"))
		d.F("//*[@id='okta-signin-password']").Input(os.Getenv("OKTA_PASS")).Input(driver.EnterKey)
		totp := gotp.NewDefaultTOTP(os.Getenv("OKTA_TOTP"))
		d.F("//*[@id='input59']").Input(totp.Now()).Input(driver.EnterKey)
	}
}

func option(opt string) func() {
	return func() {
		d.Cl("My Account")
		d.Cl(opt)
	}
}

func addAttribute(name, attType string) func() {
	return func() {
		d.Cl("Add Attribute")
		d.Cl("Name *")
		d.Keys(name)
		d.Cl("Please select...")
		d.Cl(attType)
		d.Cl("Create")
	}
}
