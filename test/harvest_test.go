package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/gost"
)

func TestHarvest(t *testing.T) {
	d, tear := gost.Gost(
		capabilities.HeadLess(),
	)
	defer tear()

	args := os.Args
	month := args[6]
	days := args[7:]
	fmt.Println(args)
	fmt.Println(month)

	d.Open(os.Getenv("HARVEST_URL"))
	d.F("Work email").Input(os.Getenv("HARVEST_USER"))
	d.F("Password").Input(os.Getenv("HARVEST_PASS"))
	d.Cl("//*[@id='log-in']")

	for _, day := range days {
		d.Cl("//*[@id='calendar-button']")
		d.Cl(fmt.Sprintf("%s %s", month, day))
		d.Cl("Copy rows from most recent timesheet")
		d.Cl("Edit entry")
		d.F("hours").Input("8:00")
		d.Cl("Update entry")
		time.Sleep(time.Second * 3)
	}

	time.Sleep(time.Second * 3)
}
