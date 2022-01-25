package test

import (
	"testing"
	"time"

	"github.com/merico-dev/lake/plugins/github/tasks"
)

func TestHandleCommitsResponse(t *testing.T) {
	ImportConfig()

	url := "https://api.github.com/repos/merico-dev/lake/commits"
	fixturePath := "fixtures/github/repos/merico-dev/lake/commits"
	resp, err := GatherDataFromApi(url, fixturePath)
	if err != nil {
		t.Fatalf("Failed to connect to api: %s", err)
	}

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	t.Fatalf("Failed to read response body: %s", err)
	// }

	// wantHeading := "joncodo"
	// bodyContent := string(body)

	// if !strings.Contains(bodyContent, wantHeading) {
	// 	t.Errorf("Heading %s not found in response", wantHeading)
	// }

	done := make(chan bool)
	// fmt.Println("JON >>> body in test", string(body))
	err = tasks.HandleCommitsResponse(resp, done)
	if err != nil {
		t.Error("Code failed to execute")
	}

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		panic("timeout")
	}
}
