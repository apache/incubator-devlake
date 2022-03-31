package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

func NewTapdApiClient(taskCtx core.TaskContext, source *models.TapdSource) (*helper.ApiAsyncClient, error) {
	// load configuration
	encKey := taskCtx.GetConfig(core.EncodeKeyEnvStr)
	auth, err := core.Decrypt(encKey, source.BasicAuthEncoded)
	if err != nil {
		return nil, fmt.Errorf("Failed to decrypt Auth Token: %w", err)
	}

	// create synchronize api client so we can calculate api rate limit dynamically
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %v", auth),
	}
	apiClient, err := helper.NewApiClient(source.Endpoint, headers, 0, source.Proxy, taskCtx.GetContext())
	if err != nil {
		return nil, err
	}
	apiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnprocessableEntity {
			return fmt.Errorf("authentication failed, please check your Basic Auth Token")
		}
		return nil
	})

	// create rate limit calculator
	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: source.RateLimit,
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

type TapdPagination struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}
