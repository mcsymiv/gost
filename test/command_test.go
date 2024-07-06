package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/command"
	"github.com/mcsymiv/gost/config"
)

func TestCmd(t *testing.T) {
	conf := config.NewConfig()
	caps := capabilities.DefaultCapabilities()

	ex, err := command.Cmd(caps, conf)
	if err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}

	defer ex.Process.Kill()
	time.Sleep(2 * time.Second)
}
