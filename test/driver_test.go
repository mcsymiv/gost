package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/command"
	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/service"
)

func start_driver(capsFn ...capabilities.CapabilitiesFunc) (*driver.WebDriver, func()) {
	d := driver.Driver()
	// setup
	return d, func() {
		// teardown
		d.Quit()
		command.OutFileLogs.Close()
		d.Command.Process.Kill()
	}
}

func setup_driver() func() {
	config.Config = config.NewConfig()

	routes := service.Handler()
	srv := service.NewServer(routes)

	go func() {
		if err := srv.Server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			panic(fmt.Sprintf("could not start server: %v", err))
		}
	}()

	return func() {
		shutdownCtx := context.Background()

		if err := srv.Server.Shutdown(shutdownCtx); err != nil {
			panic(fmt.Sprintf("HTTP shutdown error: %v", err))
		}
		fmt.Println("Graceful shutdown complete")
	}
}

func TestDriver(t *testing.T) {
	shutdown := setup_driver()
	defer shutdown()

	d, tear := start_driver(
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer tear()

	// repo := "/repository/download/"
	// allure := ":id/allure-report.zip!/allure-report-test/index.html#suites"
	host := os.Getenv("DOWNLOAD_HOST")
	// testEnv := "mt-uat"
	//
	// var rLinks []string
	// sNames := []string{
	// 	// os.Getenv("SUITE_NAME_1"), // smoke
	// 	os.Getenv("SUITE_NAME_2"), // regress
	// 	os.Getenv("SUITE_NAME_3"), // single
	// 	// os.Getenv("SUITE_NAME_4"), // m
	// 	// os.Getenv("SUITE_NAME_5"),  // ol
	// 	// os.Getenv("SUITE_NAME_6"), // hil
	// 	// os.Getenv("SUITE_NAME_7"), // gm
	// 	// os.Getenv("SUITE_NAME_8"), // business
	// 	// os.Getenv("SUITE_NAME_9"),  // visual
	// 	// os.Getenv("SUITE_NAME_10"), // iframe
	// }

	d.Url(fmt.Sprintf("%s%s", host, "/login.html"))
	time.Sleep(5 * time.Second)
	d.F("Log in using Azure Active Directory").Click()
	d.F("[aria-label^='Ending with']").Key(os.Getenv("DOWNLOAD_LOGIN")).Key(driver.EnterKey)
	// d.F("[aria-label^='Enter the password']").Key(os.Getenv("DOWNLOAD_PASS")).Key(driver.EnterKey)
	// d.F("Yes").Click()
	// d.F("Projects").Click()
	// d.F("[id='search-projects']").Key(testEnv)
	//
	// for _, sName := range sNames {
	// 	log.Println(sName)
	// 	d.F(fmt.Sprintf("//*[@data-test='sidebar']//span[contains(text(),'%s')]", sName)).Is().Click()
	// 	d.F(fmt.Sprintf("//h1//*[text()='%s']", sName)).Is()
	//
	// 	buildLinkRaw := d.F("(//*[@data-grid-root='true']//*[@data-test='ring-link'])[1]").Attr("href")
	// 	buildLink := strings.Join(strings.Split(buildLinkRaw, "/")[2:], "/")
	//
	// 	log.Println(buildLink)
	// 	rLinks = append(rLinks, fmt.Sprintf("%s%s%s%s", host, repo, buildLink, allure))
	// }
	//
	// for _, rLink := range rLinks {
	// 	d.Url(rLink)
	// 	time.Sleep(10 * time.Second)
	// 	d.F("[data-tooltip='Download CSV']").Click()
	// 	time.Sleep(10 * time.Second)
	// }
}
