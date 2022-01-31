package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

// Fire off a POST request to a live api server
// Get labels from github only in JSON body
// Poll until the "Pipelines.status == TASK_COMPLETED"

func TestMain(m *testing.M) {
	fmt.Println("***BEFORE_ALL_TESTS***")
	err := sendRequestsToLiveAPI()
	if err != nil {
		os.Exit(1)
	}
	pollForTaskCompletion()
	code := m.Run()
	os.Exit(code)
}

func pollForTaskCompletion() {
	// Its going to need a channel
	// It should block all tests until we get the right response

	// get pielines from the DB via api
	jsonStr := getPipelines()
	jsonStr.match
}

func sendRequestsToLiveAPI() error {
	getGithub()
	return nil
}

func makeAPIRequest(json []byte, url string, method string) string {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return string(body)
}

func getGithub() {
	url := "http://localhost:8080/pipelines"
	fmt.Println("URL:>", url)

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

	makeAPIRequest(jsonStr, url, "POST")
}

func getPipelines() string {
	url := "http://localhost:8080/pipelines"
	fmt.Println("URL:>", url)

	var json []byte
	jsonStr := makeAPIRequest(json, url, "GET")
	return jsonStr
}
