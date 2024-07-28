package test

import (
	"testing"
	"time"

	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/gost"
)

func TestRecord1(t *testing.T) {
	st := gost.New(t)
	defer st.Tear()

	st.Open("https://www.google.com/")

	st.Click("//*[@id='APjFqb']")

	st.Keys("hello")

	st.Keys(driver.EnterKey)

	time.Sleep(3 * time.Second)

}
