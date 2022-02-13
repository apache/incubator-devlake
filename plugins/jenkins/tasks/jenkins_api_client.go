package tasks

import (
	"net/http"

	"github.com/bndr/gojenkins"
)

type JenkinsApiClient struct {
	jenkins *gojenkins.Jenkins
}

func CreateApiClient(client *http.Client, base string, auth ...interface{}) *JenkinsApiClient {
	return &JenkinsApiClient{
		jenkins: gojenkins.CreateJenkins(client, base, auth...),
	}
}
