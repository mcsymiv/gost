package driver

import (
	"fmt"
)

// w3c Locator strategies
// src: https://www.w3.org/TR/webdriver2/#locator-strategies
const (
	ByXPath           = "xpath"
	ByLinkText        = "link text"
	ByPartialLinkText = "partial link text"
	ByTagName         = "tag name"
	ByCssSelector     = "css selector"
)

type Selector struct {
	Using, Value string
}

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
func Strategy(value string) Selector {
	var s Selector
	s.Value = value

	if value[0] == '/' || value[1] == '/' {
		s.Using = ByXPath
		return s
	}

	// if ok, m := checkSubstrings(value, ".", "#", "[", "]"); ok || m > 0 {
	if value[0] == '[' {
		s.Using = ByCssSelector
		return s
	}

	s.Using = ByXPath
	s.Value = fmt.Sprintf("//*[text()='%[1]s'] | //*[@placeholder='%[1]s'] | //*[@value='%[1]s']", value)

	return s
}

// XPathTextStrategy
// text/value based find xPathTextStrategy
// *[text()='%[1]s']
// *[@placeholder='%[1]s']
// *[@value='%[1]s']
func xPathTextStrategy(value string) Selector {
	return Selector{
		Using: ByXPath,
		Value: fmt.Sprintf("//*[text()='%[1]s'] | //*[@placeholder='%[1]s'] | //*[@value='%[1]s']", value),
	}
}

func Css(value string) Selector {
	return Selector{
		Using: ByCssSelector,
		Value: value,
	}
}

