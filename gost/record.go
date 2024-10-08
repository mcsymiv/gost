package gost

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mcsymiv/gost/config"
)

// AutoGenerated
// Chrome struct for record
type AutoGenerated struct {
	Title              string               `json:"title"`
	AutoGeneratedSteps []AutoGeneratedSteps `json:"steps"`
}

type AssertedEvents struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

// AutoGeneratedSteps
// Chrome struct of Steps record
type AutoGeneratedSteps struct {
	Type              string           `json:"type"`
	Value             string           `json:"value"`
	Key               string           `json:"key"`
	Width             int              `json:"width,omitempty"`
	Height            int              `json:"height,omitempty"`
	DeviceScaleFactor int              `json:"deviceScaleFactor,omitempty"`
	IsMobile          bool             `json:"isMobile,omitempty"`
	HasTouch          bool             `json:"hasTouch,omitempty"`
	IsLandscape       bool             `json:"isLandscape,omitempty"`
	URL               string           `json:"url,omitempty"`
	AssertedEvents    []AssertedEvents `json:"assertedEvents,omitempty"`
	Target            string           `json:"target,omitempty"`
	Selectors         [][]string       `json:"selectors,omitempty"`
	OffsetY           float32          `json:"offsetY,omitempty"`
	OffsetX           float32          `json:"offsetX,omitempty"`
}

// generatedStepSelectors
type GeneratedStepSelectors struct {
	step,
	css,
	text,
	xpath,
	aria,
	pierce string
}

type Record struct {
	AutoGenerated *AutoGenerated
	TestFile      *os.File
	RecordFile    *os.File
}

// readJsonFile
func unmarshalAutoGeneratedJson(fPath, fName string) (*AutoGenerated, error) {
	at := &AutoGenerated{}

	f, err := config.FindFile(fPath, fName)
	if err != nil {
		return nil, fmt.Errorf("error on find file: %v", err)
	}

	file, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("error on open file: %v", err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal("error on close file", err)
		}
	}()

	byteValue, _ := io.ReadAll(file)
	err = json.Unmarshal(byteValue, at)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal record.json: %v", err)
	}

	return at, nil
}

// convertSelectors
// formats chrome steps array selectors type to struct
func convertSelectors(st AutoGeneratedSteps) []*GeneratedStepSelectors {
	var genSelectors []*GeneratedStepSelectors = []*GeneratedStepSelectors{}

	if st.Type == "click" && len(st.Selectors) > 0 {
		var genSelector *GeneratedStepSelectors = &GeneratedStepSelectors{
			step: st.Type,
		}

		for _, s := range st.Selectors {
			// if strings.Contains(s[0], "aria/") {
			// 	aFormated := strings.ReplaceAll(s[0], "aria/", "")
			// 	aFormated = strings.Trim(aFormated, " ")
			// 	genSelector.aria = aFormated
			// }

			// ignore chrome selectors type if specified in record
			if strings.Contains(s[0], "pierce") {
				pFormated := strings.ReplaceAll(s[0], "pierce/", "")
				genSelector.pierce = pFormated
			}

			if strings.Contains(s[0], "xpath/") {
				xFormated := strings.ReplaceAll(s[0], "\"", "'")
				xFormated = strings.ReplaceAll(xFormated, "xpath/", "")
				genSelector.xpath = xFormated
			}

			if strings.Contains(s[0], "text/") {
				tFormated := strings.ReplaceAll(s[0], "text/", "")
				tFormated = strings.ReplaceAll(tFormated, "", "")
				tFormated = strings.ReplaceAll(tFormated, "", "")
				if tFormated != "" {
					genSelector.text = tFormated
				}
			}

			// genSelector.css = s[0]
		}

		genSelectors = append(genSelectors, genSelector)
	}

	return genSelectors
}

func writeTestStart(tName string, f *os.File) *os.File {
	var testStr string = `
package test

import (
	"testing"
	"github.com/mcsymiv/gost/gost"
)

func Test%s(t *testing.T) {
	st := gost.New(t)
	defer st.Tear()
	` // operand expected, found },

	f.WriteString(fmt.Sprintf(testStr, tName))

	return f
}

func createTestFile(fName string) (*os.File, error) {
	// creates empty test_.go file
	gF, err := os.OpenFile(fName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("error on open file: %v", err)
	}

	return gF, nil
}

// CreateSteps
// reads chrome record json
// add new steps_test.go
// with recorded selectors clicks
// with text-based selector priority
func CreateTest(fName, rName, tName string) error {
	at, err := unmarshalAutoGeneratedJson(config.Config.RecordsPath, rName)
	if err != nil {
		return fmt.Errorf("error on read record json file: %v", err)
	}

	testFile, err := createTestFile(fName)
	if err != nil {
		return fmt.Errorf("error on create test file: %v", err)
	}
	defer testFile.Close()

	testFile = writeTestStart(tName, testFile)

	var clickStr string = `
	st.Click("%s")
	` // step.Click method with %selector

	var navigateStr string = `
	st.Open("%s")
	` // step.Open method with url

	var keysStr string = `
	st.Keys("%s")
	` // step.Keys method with input text
	// and based on Active element

	var keyPressStr string = `
	st.Keys(driver.%sKey)
	` // step.Keys method with keys input

	for _, step := range at.AutoGeneratedSteps {
		if step.Type == "navigate" && !strings.Contains(step.URL, "chrome") {
			testFile.WriteString(fmt.Sprintf(navigateStr, step.URL))
		}

		if step.Type == "click" {
			genS := convertSelectors(step)
			var clickSelector string

			for _, st := range genS {
				// if st.aria != "" {
				// 	clickSelector = st.aria
				// 	break
				// }

				if st.pierce != "" {
					clickSelector = st.pierce
					break
				}

				if st.aria != "" {
					clickSelector = st.xpath
					break
				}

				if st.text != "" {
					clickSelector = st.text
					break
				}

				// TODO: upd by.Css strategy Contains logic
			}

			testFile.WriteString(fmt.Sprintf(clickStr, clickSelector))
		}

		if step.Type == "change" {
			testFile.WriteString(fmt.Sprintf(keysStr, step.Value))
		}

		if step.Type == "keyDown" {
			testFile.WriteString(fmt.Sprintf(keyPressStr, step.Key))
		}
	}

	var closeBracket string = `
}
`

	testFile.WriteString(closeBracket) // close Test%Name brakets

	return nil
}
