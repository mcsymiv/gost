package data

type DriverStatus struct {
	Message string `json:"message"`
	Ready   bool   `json:"ready"`
}

type Session struct {
	Id string `json:"sessionId"`
}

type Url struct {
	Url string `json:"url"`
}

type JsonFindUsing struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

type Selector struct {
	Using, Value string
}

// Empty
// Due to geckodriver bug: https://github.com/webdriverio/webdriverio/pull/3208
// "where Geckodriver requires POST requests to have a valid JSON body"
// Used in POST requests that don't require data to be passed by W3C
type Empty struct{}

type SendKeys struct {
	Text string `json:"text"`
}
