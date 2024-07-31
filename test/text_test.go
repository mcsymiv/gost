package test

import (
	"testing"

	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/gost"
)

func TestText(t *testing.T) {
	d, tear := gost.Gost()
	defer tear()

	d.Open("https://google.com")
	el := d.F("//*[@id='APjFqb']").Click()

	d.Keys("hello")
	d.Keys(driver.EnterKey)

	el = d.F("//*[@id='APjFqb']")
	txt := el.Text()

	if txt != "hello" {
		t.Fail()
	}

}
