package test

import (
	"fmt"
	"testing"
)

func TestStatus(t *testing.T) {
	ch := setup()
	cl := <-ch
	st, _ := cl.Status()

	fmt.Println(st)
}
