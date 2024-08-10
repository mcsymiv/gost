package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/gost/gost"
)

func TestHarvest(t *testing.T) {
	d, tear := gost.Gost()
	defer tear()

	args := os.Args
	month := args[len(args)-2]
	day := args[len(args)-1]

	d.Open(os.Getenv("HARVEST_URL"))
	d.F("Work email").Input(os.Getenv("HARVEST_USER"))
	d.F("Password").Input(os.Getenv("HARVEST_PASS"))
	d.Cl("//*[@id='log-in']")
	d.Cl("//*[@id='calendar-button']")
	d.Cl(fmt.Sprint("%s %s", month, day))
	d.Cl("Copy rows from most recent timesheet")
	d.Cl("Edit entry")
	d.F("hours").Input("8:00")
	d.Cl("Update entry")

	time.Sleep(time.Second * 3)
}
