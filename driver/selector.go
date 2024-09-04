package driver

import (
	"fmt"
	"strings"

	"github.com/mcsymiv/gost/data"
)

// checkSubstrings
// wrapper around strings.Contains
// to check multiple substrings
//
//	func checkSubstrings(str string, subs ...string) (bool, int) {
//		matches := 0
//		isCompleteMatch := true
//
//		for _, sub := range subs {
//			if strings.Contains(str, sub) {
//				matches += 1
//			} else {
//				isCompleteMatch = false
//			}
//		}
//
//		return isCompleteMatch, matches
//	}

// Strategy
// defines find element Strategy
// based on selectors "pattern"
//
// xpath:
// generally starts with forward slash /
//
// css:
// for simplicity, it check for opening bracket [
//
// text:
// as final option, if selector does not contain /, [ symbols
// XPathTextStrategy will be used
func Strategy(value string) *data.Selector {
	var s *data.Selector = &data.Selector{}
	s.Value = value

	if value[0] == '/' || value[1] == '/' {
		s.Using = data.ByXPath
		return s
	}

	// if ok, m := checkSubstrings(value, ".", "#", "[", "]"); ok || m > 0 {
	if value[0] == '[' || value[0] == '#' {
		s.Using = data.ByCssSelector
		return s
	}

	txt := []string{
		"//*[text()='%[1]s']",
		"//*[@placeholder='%[1]s']",
		"//*[@value='%[1]s']",
		"//*[@title='%[1]s']",
		"//*[@aria-label='%[1]s']",
	}
	xpathText := strings.Join(txt, " | ")

	s.Using = data.ByXPath
	s.Value = fmt.Sprintf(xpathText, value)

	return s
}

func NextXpathStrategy(value string) *data.Selector {
	var s *data.Selector = &data.Selector{}

	txt := []string{
		".//*[text()='%[1]s']",
		".//*[@placeholder='%[1]s']",
		".//*[@value='%[1]s']",
		".//*[@title='%[1]s']",
		".//*[@aria-label='%[1]s']",
	}
	xpathText := strings.Join(txt, " | ")

	s.Using = data.ByXPath
	s.Value = fmt.Sprintf(xpathText, value)

	return s
}

// XPathTextStrategy
// text/value based find xPathTextStrategy
// *[text()='%[1]s']
// *[@placeholder='%[1]s']
// *[@value='%[1]s']
func xPathTextStrategy(value string) *data.Selector {
	return &data.Selector{
		Using: data.ByXPath,
		Value: fmt.Sprintf("//*[text()='%[1]s'] | //*[@placeholder='%[1]s'] | //*[@value='%[1]s']", value),
	}
}

func Css(value string) *data.Selector {
	return &data.Selector{
		Using: data.ByCssSelector,
		Value: value,
	}
}
