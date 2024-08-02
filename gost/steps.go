package gost

import (
	"fmt"
	"testing"
	"time"

	"github.com/mcsymiv/gost/capabilities"
	"github.com/mcsymiv/gost/config"
	"github.com/mcsymiv/gost/driver"
)

type Step struct {
	TK     *testing.T
	WD     *driver.WebDriver
	Tear   func()
	Config config.WebConfig
}

func New(t *testing.T, capsFn ...capabilities.CapabilitiesFunc) *Step {
	wd, tear := Gost(capsFn...)
	return &Step{
		TK:   t,
		WD:   wd,
		Tear: tear,
	}
}

func (s *Step) Open(url string) {
	open := func() error {
		_, err := s.WD.WebClient.Open(url, s.WD.SessionId)
		if err != nil {
			s.WD.Screenshot()
			return fmt.Errorf("error on open")
		}

		return nil
	}

	err := open()
	if err != nil {
		s.TK.Errorf("%v", err)
	}
}

func (s *Step) Click(selector string) {
	find := func() (*driver.WebElement, error) {
		selector := driver.Strategy(selector)
		eId, err := s.WD.WebClient.FindElement(selector, s.WD.SessionId)
		if err != nil {
			s.WD.Screenshot()
			return nil, err
		}

		return &driver.WebElement{
			WebDriver:    s.WD,
			WebElementId: eId,
		}, nil
	}

	el, err := find()
	if err != nil {
		s.TK.Error(err)
	}

	click := func() (*driver.WebElement, error) {
		err := s.WD.WebClient.Click(s.WD.SessionId, el.WebElementId)
		if err != nil {
			s.WD.Screenshot()
			return nil, fmt.Errorf("error on click: %v", err)
		}

		return el, nil
	}

	_, err = click()
	if err != nil {
		s.TK.Errorf("%v", err)
	}
}

func (s *Step) TryClick(selectors ...string) {
	find := func(sel string) (*driver.WebElement, error) {
		selector := driver.Strategy(sel)
		eId, err := s.WD.WebClient.FindElement(selector, s.WD.SessionId)
		if err != nil {
			s.WD.Screenshot()
			return nil, fmt.Errorf("error on find element: %v", err)
		}

		return &driver.WebElement{
			WebDriver:    s.WD,
			WebElementId: eId,
		}, nil
	}

	var el *driver.WebElement
	var err error
	for _, selector := range selectors {
		el, err = find(selector)
		if err != nil {
			continue
		}

		break
	}

	if el == nil || err != nil {
		s.TK.Errorf("%v", err)
	}

	click := func() (*driver.WebElement, error) {
		err := s.WD.WebClient.Click(s.WD.SessionId, el.WebElementId)
		if err != nil {
			s.WD.Screenshot()
			return nil, fmt.Errorf("error on click: %v", err)
		}

		return el, nil
	}

	_, err = click()
	if err != nil {
		s.TK.Errorf("%v", err)
	}
}

// Type
// Sends keys onto active element
// after click
func (s *Step) Input(text, selector string) {
	find := func() (*driver.WebElement, error) {
		selector := driver.Strategy(selector)
		eId, err := s.WD.WebClient.FindElement(selector, s.WD.SessionId)
		if err != nil {
			s.WD.Screenshot()
			return nil, fmt.Errorf("error on find element: %v", err)
		}

		return &driver.WebElement{
			WebDriver:    s.WD,
			WebElementId: eId,
		}, nil
	}

	el, err := find()
	if err != nil {
		s.TK.Errorf("%v", err)
	}

	click := func() (*driver.WebElement, error) {
		err := s.WD.WebClient.Click(s.WD.SessionId, el.WebElementId)
		if err != nil {
			s.WD.Screenshot()
			return nil, fmt.Errorf("error on click: %v", err)
		}

		return el, nil
	}

	el, err = click()
	if err != nil {
		s.TK.Errorf("%v", err)
	}

	keys := func() (*driver.WebElement, error) {
		err := s.WD.WebClient.Input(text, s.WD.SessionId, el.WebElementId)
		if err != nil {
			s.WD.Screenshot()
			return nil, fmt.Errorf("error on keys: %v", err)
		}

		return el, nil
	}

	_, err = keys()
	if err != nil {
		s.TK.Errorf("%v", err)
	}
}

func (s *Step) Keys(text string) {
	action := func() error {
		err := s.WD.WebClient.Action(text, string(driver.KeyDownAction), s.WD.SessionId)
		if err != nil {
			s.WD.Screenshot()
			return fmt.Errorf("error on find element: %v", err)
		}

		return nil
	}

	err := action()
	if err != nil {
		s.TK.Errorf("%v", err)
	}
}

func (s *Step) Is(selector string) bool {
	find := func() (*driver.WebElement, error) {
		selector := driver.Strategy(selector)
		eId, err := s.WD.WebClient.FindElement(selector, s.WD.SessionId)
		if err != nil {
			s.WD.Screenshot()
			return nil, fmt.Errorf("error on find element: %v", err)
		}

		return &driver.WebElement{
			WebDriver:    s.WD,
			WebElementId: eId,
		}, nil
	}

	el, err := find()
	if err != nil {
		s.TK.Errorf("%v", err)
	}

	is := func() (bool, error) {
		ok, err := s.WD.WebClient.Is(el.WebElementId, s.WD.SessionId)
		if err != nil {
			s.WD.Screenshot()
			return false, fmt.Errorf("error on find element: %v", err)
		}

		return ok, nil
	}

	ok, err := is()
	if err != nil {
		s.TK.Errorf("%v", err)
	}

	return ok
}

func (s *Step) Until(fn func() bool) {
	var success bool
	start := time.Now()
	end := start.Add(s.Config.WaitForTimeout * time.Second)

	for {
		success = fn()
		if success {
			break
		}

		if time.Now().After(end) {
			panic("error on until")
		}

		time.Sleep(s.Config.WaitForInterval * time.Millisecond)
	}
}
