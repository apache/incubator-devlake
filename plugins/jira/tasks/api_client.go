package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

func NewJiraApiClient(taskCtx core.TaskContext, connection *models.JiraConnection) (*helper.ApiAsyncClient, error) {
	// load configuration
	encKey := taskCtx.GetConfig(core.EncodeKeyEnvStr)
	auth, err := core.Decrypt(encKey, connection.BasicAuthEncoded)
	if err != nil {
		return nil, fmt.Errorf("Failed to decrypt Auth Token: %w", err)
	}

	// create synchronize api client so we can calculate api rate limit dynamically
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %v", auth),
	}
	apiClient, err := helper.NewApiClient(connection.Endpoint, headers, 0, connection.Proxy, taskCtx.GetContext())
	if err != nil {
		return nil, err
	}
	apiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your Basic Auth Token")
		}
		return nil
	})

	// create rate limit calculator
	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: connection.RateLimit,
	}
	asyncApiClient, err := helper.CreateAsyncApiClient(
		taskCtx,
		apiClient,
		rateLimiter,
	)
	if err != nil {
		return nil, err
	}

	return asyncApiClient, nil
}

type JiraPagination struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

func GetJiraServerInfo(client *helper.ApiAsyncClient) (*models.JiraServerInfo, int, error) {
	res, err := client.Get("api/2/serverInfo", nil, nil)
	if err != nil {
		return nil, 0, err
	}
	if res.StatusCode >= 300 || res.StatusCode < 200 {
		return nil, res.StatusCode, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}
	serverInfo := &models.JiraServerInfo{}
	err = helper.UnmarshalResponse(res, serverInfo)
	if err != nil {
		return nil, res.StatusCode, err
	}
	return serverInfo, res.StatusCode, nil
}
