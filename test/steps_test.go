package test

import (
	"testing"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/gost"
)

func TestSteps(t *testing.T) {
	d, tear, shutdown, st := gost.Gost(
		t,
		capabilities.MozPrefs("intl.accept_languages", "en-GB"),
	)
	defer shutdown()
	defer tear()

	st.It("adding functions",
		d.FindClick(""),
		d.FindClick(""),
	).Run()

}
