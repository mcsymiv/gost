package gost

import (
	"fmt"
	"testing"

	"github.com/mcsymiv/gost/driver"
)

type Step struct {
	TK   *testing.T
	WD   *driver.WebDriver
	Tear func()
}

func New(t *testing.T) *Step {
	wd, tear := Gost()
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

	_, err = click()
	if err != nil {
		s.TK.Errorf("%v", err)
	}
}

func (s *Step) Type(text, selector string) {
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
		err := s.WD.WebClient.Keys(text, s.WD.SessionId, el.WebElementId)
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
