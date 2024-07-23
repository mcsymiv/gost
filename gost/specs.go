package gost

import (
	"fmt"
	"testing"

	"github.com/mcsymiv/gost/driver"
	"github.com/mcsymiv/gost/service"
)

type Fn func() *Spec

type Spec struct {
	TK   *testing.T
	F    []Fn
	WD   *driver.WebDriver
	SRV  *service.WebServer
	Tear []func()
}

func (s *Spec) Shutdown() {
	for _, it := range s.Tear {
		it()
	}
}

// func (s *Spec) Run(desc string, fn ...Fn) {
// 	var fns []Fn
// 	for _, it := range fn {
// 		fns = append(fns, it)
// 	}
//
// 	s.F = fns
//
// 	var res *Spec
// 	for _, it := range s.F {
// 		res = it()
// 		if res == nil {
// 			s.TK.Error()
// 		}
// 		break
// 	}
//
// 	s.TK.Error("failed test")
// }

func (sp *Spec) Click(s string) *Spec {
	findElement := func() (*driver.WebElement, error) {
		selector := driver.Strategy(s)
		eId, err := sp.WD.WebClient.FindElement(selector, sp.WD.SessionId)
		if err != nil {
			return nil, fmt.Errorf("error on find element: %v", err)
		}

		return &driver.WebElement{
			WebDriver:    sp.WD,
			WebElementId: eId,
		}, nil
	}

	el, err := findElement()
	if err != nil {
		// send stop to spec with error
	}

	clickElement := func() (*driver.WebElement, error) {
		err := sp.WD.WebClient.Click(sp.WD.SessionId, el.WebElementId)
		if err != nil {
			return nil, fmt.Errorf("error on click: %v", err)
		}

		return el, nil
	}

	_, err = clickElement()
	if err != nil {
		// send stop to spec with error
	}

	return sp
}

func (sp *Spec) Open(s string) *Spec {
	open := func() error {
		_, err := sp.WD.WebClient.Open(s, sp.WD.SessionId)
		if err != nil {
			return fmt.Errorf("error on open: %v", err)
		}

		return nil
	}

	err := open()
	if err != nil {
		// send stop to spec with error
	}

	return sp
}
