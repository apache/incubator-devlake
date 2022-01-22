package tasks

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestHandleCommitsResponse(t *testing.T) {
	ImportConfig()

	url := "https://api.github.com/repos/merico-dev/lake/commits"
	fixturePath := "fixtures/github/repos/merico-dev/lake/commits"
	resp, err := GatherDataFromApi(url, fixturePath)
	if err != nil {
		t.Fatalf("Failed to connect to api: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	wantHeading := "joncodo"
	bodyContent := string(body)

	if !strings.Contains(bodyContent, wantHeading) {
		t.Errorf("Heading %s not found in response", wantHeading)
	}
}
