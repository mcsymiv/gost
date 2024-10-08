package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mcsymiv/gost/gost"
)

func TestHome(t *testing.T) {
	d, tear := gost.Gost()
	defer tear()

	fmt.Println(os.Args)

	d.Url("http://192.168.0.1/")
	d.F("//*[@type='password']").Input(os.Getenv("HOME_PASS"))
	time.Sleep(1 * time.Second)
	d.Cl("LOG IN")
	time.Sleep(3 * time.Second)
	d.Cl("Clients")
	time.Sleep(7 * time.Second)
}
