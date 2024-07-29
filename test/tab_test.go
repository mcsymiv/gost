package test

import (
	"testing"
	"time"

	"github.com/mcsymiv/gost/gost"
)

func TestTab(t *testing.T) {
	d, tear := gost.Gost()

	defer tear()

	d.Open("https://google.com")
	d.NewTab()

	tabs := d.Tabs()
	if len(tabs) != 2 {
		t.Fail()
	}

	d.Tab(1)

	time.Sleep(time.Second * 4)

	d.Keys()

}
