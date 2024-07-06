package command

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/config"
)

var GeckoDriverPath string = "geckodriver"
var ChromeDriverPath string = "chromedriver"

func Cmd(caps *capabilities.Capabilities, conf *config.WebConfig) (*exec.Cmd, error) {
	// returns command arguments for specified driver to start from shell
	var cmdArgs []string = driverCommand(caps)

	// previously used line to start driver
	// cmd := exec.Command("zsh", "-c", GeckoDriverrequest, "--port", "4444", ">", "logs/gecko.session.logs", "2>&1", "&")
	// open the out file for writing
	OutFileLogs, err := os.Create(os.Getenv("DRIVER_LOGS"))
	if err != nil {
		panic(fmt.Sprintf("failed to start driver service: %v", err))
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = OutFileLogs
	cmd.Stderr = OutFileLogs

	err = cmd.Start()
	if err != nil {
		panic(fmt.Sprintf("failed to start driver service: %v", err))
	}

	// delay to wait for driver service to start
	// change to status polling
	time.Sleep(1 * time.Second)

	if cmd.Process.Pid == 0 {
		panic(fmt.Sprintf("failed to start driver service: %v", err))
	}

	return cmd, nil
}

// driverCommand
// Check for specified driver/browser name to pass to cmd to start the driver server
func driverCommand(cap *capabilities.Capabilities) []string {
	// when calling /bin/zsh -c command
	// command arguments will be ignored
	var cmdArgs []string = []string{
		// "-c",
	}

	if cap.Capabilities.AlwaysMatch.BrowserName == "firefox" {
		cmdArgs = append(cmdArgs, GeckoDriverPath, "--port", "4444", "--log", "trace")
	} else {
		cmdArgs = append(cmdArgs, ChromeDriverPath, fmt.Sprintf("--port=%s", "4444"))
	}

	// redirect output argumetns ignored when used in exec.Command
	// cmdArgs = append(cmdArgs, ">", "logs/session.log", "2>&1", "&")
	return cmdArgs
}
