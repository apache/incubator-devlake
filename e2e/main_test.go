package e2e

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/merico-dev/lake/plugins/core"
)

type PipelinesAPIResponse struct {
	Pipelines []struct {
		ID     int
		Status string
	}
}

var Db sql.DB

func TestMain(m *testing.M) {
	if os.Getenv("E2E") != "true" {
		return
	}
	fmt.Println("***BEFORE_ALL_TESTS***")
	err := sendRequestsToLiveAPI()
	if err != nil {
		panic(err)
	}

	loopDelay := 1
	readyToTest := false

	// Block all test from running until all tasks are done in the pipelines table
	// Poll until there are no pipelines with pipeline.Status == "TASK_RUNNING" || pipeline.Status == "TASK_CREATED"
	for !readyToTest {
		readyToTest, err = checkForTaskCompletion()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Duration(loopDelay * int(time.Second)))
	}
	code := m.Run()
	os.Exit(code)
}

func checkForTaskCompletion() (bool, error) {
	// get pipelines from the DB via api
	pipelines, err := getPipelines()
	if err != nil {
		return false, err
	}
	for _, pipeline := range pipelines.Pipelines {
		// make sure all tasks are done
		if pipeline.Status == "TASK_RUNNING" || pipeline.Status == "TASK_CREATED" {
			return false, nil
		}
	}
	return true, nil
}

// Send off all requests to the api to gather all data before we run our tests
func sendRequestsToLiveAPI() error {
	// getGithub()
	getGitLab()
	return nil
}

func makeAPIRequest(json []byte, url string, method string, v interface{}) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	err = core.UnmarshalResponse(resp, v)
	if err != nil {
		return err
	}
	return nil
}

func getGitLab() error {
	url := "http://localhost:8080/pipelines"

	var jsonStr = []byte(`{
        "name": "test-all",
        "tasks": [
            [
                {
                    "Plugin": "gitlab",
                    "Options": {
											"projectId": 28270340,
											"tasks": ["collectMrs"]
                    }
                }
            ]
        ]
    }`)

	err := makeAPIRequest(jsonStr, url, "POST", nil)
	if err != nil {
		return err
	}
	return nil
}

// Gather all data from the github plugin

// "collectRepo" is not valid.
func getGithub() error {
	url := "http://localhost:8080/pipelines"

	var jsonStr = []byte(`{
        "name": "test-all",
        "tasks": [
            [
                {
                    "Plugin": "github",
                    "Options": {
                        "repositoryName": "lake",
                        "owner": "merico-dev",
                        "tasks": ["collectRepo"]
                    }
                }
            ]
        ]
    }`)

	err := makeAPIRequest(jsonStr, url, "POST", nil)
	if err != nil {
		return err
	}
	return nil
}

// Get the list of all pipelines so we can see when collection is done
func getPipelines() (*PipelinesAPIResponse, error) {
	url := "http://localhost:8080/pipelines"

	var json []byte
	pipelinesAPIResponse := &PipelinesAPIResponse{}
	err := makeAPIRequest(json, url, "GET", pipelinesAPIResponse)
	if err != nil {
		return nil, err
	}
	return pipelinesAPIResponse, nil
}
