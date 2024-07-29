package test

import (
	"testing"

	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/gost"
)

func TestRecord(t *testing.T) {
	config.Config = config.NewConfig()

	var fName string = "../test/rec_2_test.go"
	var rName string = "rec_2.json"
	var tName string = "Record2"

	gost.CreateTest(fName, rName, tName)
}
