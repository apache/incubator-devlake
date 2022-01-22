package tasks

import (
	"net/http"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/merico-dev/lake/config"
)

func ImportConfig() {
	_ = config.LoadConfigFile("../../../.env")
}

func GatherDataFromApi(url string, path string) (*http.Response, error) {
	// Start our recorder
	r, err := recorder.New(path)
	if err != nil {
		return nil, err
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	// Create an HTTP client and inject our transport
	client := &http.Client{
		Transport: r, // Inject as transport!
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
