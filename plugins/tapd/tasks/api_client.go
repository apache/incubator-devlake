package tasks

import (
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

func NewTapdApiClient(taskCtx core.TaskContext, connection *models.TapdConnection) (*helper.ApiAsyncClient, error) {
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
	apiClient, err := helper.NewApiClient(connection.Endpoint, headers, 0, "", taskCtx.GetContext())
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
