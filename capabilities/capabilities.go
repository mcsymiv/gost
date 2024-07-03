package capabilities

type Capabilities struct {
	Capabilities BrowserCapabilities `json:"capabilities"`
}

type BrowserCapabilities struct {
	AlwaysMatch `json:"alwaysMatch"`
}

type AlwaysMatch struct {
	AcceptInsecureCerts bool   `json:"acceptInsecureCerts"`
	BrowserName         string `json:"browserName"`
	Timeouts            `json:"timeouts,omitempty"`
	ChromeOptions       `json:"goog:chromeOptions,omitempty"`
	MozOptions          `json:"moz:firefoxOptions,omitempty"`
	PageLoad            string `json:"pageLoadStrategy,omitempty"`
}

type Timeouts struct {
	Implicit float32 `json:"implicit,omitempty"`
	Script   float32 `json:"script,omitempty"`
}

type ChromeOptions struct {
	Binary string   `json:"binary,omitempty"`
	Args   []string `json:"args,omitempty"`
}

type MozOptions struct {
	Profile string            `json:"profile,omitempty"`
	Binary  string            `json:"binary,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Prefs   map[string]string `json:"prefs,omitempty"`
	Log     `json:"log,omitempty"`
}

type Log struct {
	Level string `json:"level,omitempty"`
}

type CapabilitiesFunc func(*Capabilities)

// DefaultCapabilities
// Sets default firefox browser with local dev url
// With defined in service port, i.e. :4444
// Port and Host fields are used and passed to the WebDriver instance
// To reference and build current driver url
func DefaultCapabilities() *Capabilities {
	return &Capabilities{
		Capabilities: BrowserCapabilities{
			AlwaysMatch{
				AcceptInsecureCerts: true,
				BrowserName:         "chrome",
				PageLoad:            "eager",
				Timeouts: Timeouts{
					Implicit: 1000,
				},
			},
		},
	}
}

func ImplicitWait(w float32) CapabilitiesFunc {
	return func(cap *Capabilities) {
		cap.Capabilities.AlwaysMatch.Timeouts.Implicit = w
	}
}

// PageLoadStrategy
// https://html.spec.whatwg.org/#current-document-readiness
func PageLoadStrategy(st string) CapabilitiesFunc {
	return func(cap *Capabilities) {
		cap.Capabilities.AlwaysMatch.PageLoad = st
	}
}

func HeadLess() CapabilitiesFunc {
	return func(cap *Capabilities) {
		cap.Capabilities.AlwaysMatch.MozOptions = MozOptions{
			Args: []string{"-headless"},
		}
	}
}

func ChromeArgs(args []string) CapabilitiesFunc {
	return func(cap *Capabilities) {
		cap.Capabilities.AlwaysMatch.ChromeOptions = ChromeOptions{
			Args: args,
		}
	}
}

func BrowserName(b string) CapabilitiesFunc {
	return func(cap *Capabilities) {
		cap.Capabilities.AlwaysMatch.BrowserName = b
	}
}

func MozPrefs(k, v string) CapabilitiesFunc {
	return func(caps *Capabilities) {
		caps.Capabilities.MozOptions.Prefs = map[string]string{k: v}
	}
}
