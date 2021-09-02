package tasks

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
)

var apiClient *core.ApiClient

func GetJiraApiClient() *core.ApiClient {
	if apiClient == nil {
		apiClient = &core.ApiClient{}
		apiClient.Setup(
			config.V.GetString("JIRA_ENDPOINT"),
			map[string]string{
				"Authorization": fmt.Sprintf("Basic %v", config.V.GetString("JIRA_BASIC_AUTH_ENCODED")),
			},
			10*time.Second,
			3,
		)
	}
	return apiClient
}
