package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

const (
	// LegacyWebElementIdentifier is the string constant used in the old Selenium 2 protocol
	// WebDriver JSON protocol that is the key for the map that contains an
	// unique element identifier.
	// This value is ignored in element id retreival
	LegacyWebElementIdentifier = "ELEMENT"

	// WebElementIdentifier is the string constant defined by the W3C Selenium 3 protocol
	// specification that is the key for the map that contains a unique element identifier.
	WebElementIdentifier = "element-6066-11e4-a52e-4f735466cecf"

	// ShadowRootIdentifier A shadow root is an abstraction used to identify a shadow root when
	// it is transported via the protocol, between remote and local ends.
	ShadowRootIdentifier = "shadow-6066-11e4-a52e-4f735466cecf"
)

var Config *WebConfig

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../..", "gost")
)

const ApplicationJson string = "application/json"
const ContenType string = "Content-Type"

type WebConfig struct {
	// WebServerAddr
	// Default value http://localhost:8080
	WebServerPort string
	WebServerAddr string

	// WebDriverAddr
	// Default value http://localhost:4444
	WebDriverPort string
	WebDriverAddr string

	// DriverLogsFile
	DriverLogsFile string

	// ConfigFile
	// name of cinfiguration file
	ConfigFile string

	// ScreenshotOnFail
	// used in find element strategy
	// takes screenshot and writes to artifacts
	// if unable to find webelement within timeout
	ScreenshotOnFail bool

	// WaitForTimeout
	// used in find element strategy
	// controls timeout of performing driver.F("selector") find
	// 20 seconds default value
	WaitForTimeout time.Duration

	// WaitForInterval
	// delay to retry find element request
	// 200 ms is an arbitrary value
	WaitForInterval time.Duration

	// RefreshOnFindError
	// calls /session/{sessionId}/refresh
	// if find retry fails
	RefreshOnFindError bool

	// Artifact path
	//
	// ArtifactRecordsPath
	// from app root a directory that stores
	// Google Chrome Recorder json files
	// for TestSteps generation
	// will check specified path for *.json records
	RecordsPath string

	// ArtifactScreenshotsPath
	// from app root a directory where
	// ScreenshotOnFail, or d.Screenshot()
	// stores driver screnshots in *.jpg format
	ScreenshotsPath string

	// ArtifactJsFilesPath
	// from app root a directory where
	JsFilesPath string
}

type ConfigFunc func(*WebConfig)

func DefaultConfig() *WebConfig {
	return &WebConfig{
		WebServerAddr:    "http://localhost:8080",
		WebDriverAddr:    "http://localhost:4444",
		DriverLogsFile:   GetRoot("driver.logs"),
		ConfigFile:       ".config",
		ScreenshotOnFail: true,
		WaitForTimeout:   20,
		WaitForInterval:  200,
		JsFilesPath:      GetPath("js"),
		ScreenshotsPath:  GetPath("screenshots"),
		RecordsPath:      GetPath("records"),
	}
}

func NewConfig(confFn ...ConfigFunc) *WebConfig {
	var conf *WebConfig

	err := loadEnv(Root, ".config")
	if err != nil {
		conf = DefaultConfig()

		return conf
	}

	conf = &WebConfig{
		WebServerPort:    os.Getenv("SERVER_PORT"),
		WebServerAddr:    fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")),
		WebDriverPort:    os.Getenv("DRIVER_PORT"),
		WebDriverAddr:    fmt.Sprintf("%s:%s", os.Getenv("DRIVER_HOST"), os.Getenv("DRIVER_PORT")),
		DriverLogsFile:   GetRoot(os.Getenv("DRIVER_LOGS")),
		WaitForTimeout:   toWaitTimeout(os.Getenv("WAIT_TIMEOUT")),
		WaitForInterval:  toWaitInterval(os.Getenv("WAIT_INTERVAL")),
		JsFilesPath:      GetPath(os.Getenv("JS_FILES_PATH")),
		ScreenshotsPath:  GetPath(os.Getenv("SCREENSHOTS_PATH")),
		RecordsPath:      GetPath(os.Getenv("RECORDS_PATH")),
		ScreenshotOnFail: WebConfigDriverScreenshoOnFail(os.Getenv("SCREENSHOT_ON_FAIL")),
	}

	return conf
}

func WebConfigServerAddr(addr string) ConfigFunc {
	return func(conf *WebConfig) {
		conf.WebServerAddr = addr
	}
}

func WebConfigDriverAddr(addr string) ConfigFunc {
	return func(conf *WebConfig) {
		conf.WebDriverAddr = addr
	}
}

func WebConfigDriverScreenshoOnFail(onFail string) bool {
	f, err := strconv.ParseBool(onFail)
	if err != nil {
		return f
	}

	return true
}

func toWaitTimeout(dur string) time.Duration {
	d, err := strconv.Atoi(dur)
	if err != nil {
		return 20
	}

	return time.Duration(d)
}

func toWaitInterval(dur string) time.Duration {
	d, err := strconv.Atoi(dur)
	if err != nil {
		return 200
	}

	return time.Duration(d)
}
