package test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/gost"
)

func findSuite(suite string) string {
	var suites []string
	var env string

	for i := 1; i < 12; i++ {
		suites = append(suites, os.Getenv(fmt.Sprintf("SUITE_NAME_%s", strconv.Itoa(i))))
	}

	for _, s := range suites {
		name := strings.ToLower(s)

		if strings.Contains(name, suite) {
			env = s
		}

	}

	return env
}

func trigger(action, env, suite string) {
	d, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	host := os.Getenv("TC_HOST")
	d.Url(fmt.Sprintf("%s%s", host, "/login.html"))
	d.Cl("Log in using Azure Active Directory")
	d.F("[aria-label^='Ending with']").Input(os.Getenv("TC_LOGIN")).Input(driver.EnterKey)
	d.F("[aria-label^='Enter the password']").Input(os.Getenv("TC_PASS")).Input(driver.EnterKey)
	d.Cl("Yes")
	d.F("Projects").Click()
	d.F("[id='search-projects']").Input(env)

	sName := findSuite(suite)
	d.F(fmt.Sprintf("//*[@data-test='sidebar']//span[contains(text(),'%s')]", sName)).Click()

	d.F("Edit configuration...").Is().Click()
	d.Cl("Triggers")
	d.Cl("[id*='triggerActionsTRIGGER']")

	if action == "disable" {
		d.Cl("Disable trigger")
		d.Cl("Disable")
	}

	d.Cl("Enable trigger")

	time.Sleep(time.Second * 3)
}

func run(env, suite string) {
	d, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	host := os.Getenv("TC_HOST")
	d.Url(fmt.Sprintf("%s%s", host, "/login.html"))
	d.Cl("Log in using Azure Active Directory")
	d.F("[aria-label^='Ending with']").Input(os.Getenv("TC_LOGIN")).Input(driver.EnterKey)
	d.F("[aria-label^='Enter the password']").Input(os.Getenv("TC_PASS")).Input(driver.EnterKey)
	d.Cl("Yes")
	d.F("Projects").Click()
	d.F("[id='search-projects']").Input(env)

	sName := findSuite(suite)
	d.F(fmt.Sprintf("//*[@data-test='sidebar']//span[contains(text(),'%s')]", sName)).Is().Click()
	// d.F(fmt.Sprintf("//h1//*[text()='%s']", sName)).IsDisplayed()

	d.F("Run").IsDisplayed()
	d.Cl("Run")

	time.Sleep(time.Second * 3)
}

func result(testEnv string) {
	d, tear := gost.Gost(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	repo := "/repository/download/"
	allure := ":id/allure-report.zip!/allure-report-test/index.html#suites"
	host := os.Getenv("TC_HOST")

	var rLinks []string
	sNames := []string{
		// os.Getenv("SUITE_NAME_1"), // smoke
		os.Getenv("SUITE_NAME_2"), // regress
		os.Getenv("SUITE_NAME_3"), // single
		// os.Getenv("SUITE_NAME_4"), // m
		// os.Getenv("SUITE_NAME_5"),  // ol
		// os.Getenv("SUITE_NAME_6"), // hil
		// os.Getenv("SUITE_NAME_7"), // gm
		// os.Getenv("SUITE_NAME_8"), // business
		// os.Getenv("SUITE_NAME_9"),  // visual
		// os.Getenv("SUITE_NAME_10"), // iframe
	}

	d.Url(fmt.Sprintf("%s%s", host, "/login.html"))
	// time.Sleep(5 * time.Second)
	d.F("Log in using Azure Active Directory").Click()
	d.F("[aria-label^='Ending with']").Input(os.Getenv("TC_LOGIN")).Input(driver.EnterKey)
	d.F("[aria-label^='Enter the password']").Input(os.Getenv("TC_PASS")).Input(driver.EnterKey)
	d.F("Yes").Click()
	d.F("Projects").Click()
	d.F("[id='search-projects']").Input(testEnv)

	for _, sName := range sNames {
		fmt.Println(sName)
		d.F(fmt.Sprintf("//*[@data-test='sidebar']//span[contains(text(),'%s')]", sName)).Is().Click()
		d.F(fmt.Sprintf("//h1//*[text()='%s']", sName)).IsDisplayed()

		buildLinkRaw := d.F("(//*[@data-grid-root='true']//*[@data-test='ring-link'])[1]").Attr("href")
		buildLink := strings.Join(strings.Split(buildLinkRaw, "/")[2:], "/")

		fmt.Println(buildLink)
		rLinks = append(rLinks, fmt.Sprintf("%s%s%s%s", host, repo, buildLink, allure))
	}

	for _, rLink := range rLinks {
		d.Url(rLink)
		time.Sleep(10 * time.Second)
		d.F("[data-tooltip='Download CSV']").Click()
		time.Sleep(10 * time.Second)
	}
}

func TestTc(t *testing.T) {
	args := os.Args[6:]

	if len(args) != 0 {
		switch args[0] {
		case "trigger":
			trigger(args[1], args[2], args[3])
		case "result":
			result(args[1])
		case "run":
			run("dev01", "api")
		default:
			fmt.Println(`
				Usage:
				make tc trigger enable dev01 smoke
				make tc result dev01
			`)
		}
	}
}
