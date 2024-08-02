package test

import (
	"testing"

	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/gost"
)

func TestRecord(t *testing.T) {
	config.Config = config.NewConfig()

	var fName string = "../test/rec_3_test.go"
	var rName string = "rec_3.json"
	var tName string = "Record3"

	gost.CreateTest(fName, rName, tName)
}
