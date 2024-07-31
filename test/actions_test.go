package test

import (
	"testing"
	"time"

	"github.com/mcsymiv/gost/gost"
)

func TestAction(t *testing.T) {
	d, tear := gost.Gost()
	defer tear()

	d.Open("https://google.com")
	d.F("//*[@id='APjFqb']").Click()

	d.Keys("hello")

	time.Sleep(time.Second * 3)
}
