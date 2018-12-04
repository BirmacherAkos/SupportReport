package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bitrise-io/go-utils/log"
	prompt "github.com/c-bata/go-prompt"
)

var completerOptions []prompt.Suggest

func completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(completerOptions, d.GetWordBeforeCursor(), true)
}

// BasicData ...
type BasicData struct {
	FormatVersion        string `json:"format_version"`
	GeneratedAtTimestamp int    `json:"generated_at_timestamp"`
	SteplibSource        string `json:"steplib_source"`
	DownloadLocations    []struct {
		Type string `json:"type"`
		Src  string `json:"src"`
	} `json:"download_locations"`
	AssetsDownloadBaseURI string          `json:"assets_download_base_uri"`
	Steps                 map[string]Step `json:"steps"`
}

// Step ...
type Step struct {
	Info struct {
		AssetUrls struct {
			IconSvg string `json:"icon.svg"`
		} `json:"asset_urls"`
	} `json:"info"`
	LatestVersionNumber string                            `json:"latest_version_number"`
	Versions            map[string]map[string]interface{} `json:"versions"`
	name                string
}

func fetchSteps() (BasicData, error) {
	response, err := http.Get("https://bitrise-steplib-collection.s3.amazonaws.com/spec.json")
	if err != nil {
		return BasicData{}, err
	}

	var data BasicData

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return BasicData{}, err
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return BasicData{}, err
	}

	return data, nil
}

func logPretty(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	return fmt.Sprintf("%+v\n", string(b))
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func main() {
	log.Infof("Fetching step list")

	d, err := fetchSteps()
	if err != nil {
		failf("Failed to fetch the step list from the server, error: %s", err)
	}

	var names []string
	for stepName, step := range d.Steps {
		step.name = stepName
		names = append(names, stepName)
	}

	log.Printf(logPretty(names))

	// question := `Which step failed?`
	// failingStepName, err := goinp.SelectFromStrings(question, names)

	// question = fmt.Sprintf("Which version failed of the step (%s)?", failingStepName)

	// keys := reflect.ValueOf(d.Steps[failingStepName].Versions).MapKeys()
	// stepVersions := make([]string, len(keys))
	// for i := 0; i < len(keys); i++ {
	// 	stepVersions[i] = keys[i].String()
	// }
	// _, err = goinp.SelectFromStrings(question, stepVersions)

	for _, name := range names {
		completerOptions = append(completerOptions, prompt.Suggest{Text: name, Description: "No descp yest"})
	}

	fmt.Println("Please select table.")
	t := prompt.Input("> ", completer)
	fmt.Println("You selected " + t)

}
