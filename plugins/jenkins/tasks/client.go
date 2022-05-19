package tasks

import (
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/utils"
)

func CreateApiClient(taskCtx core.TaskContext) (*helper.ApiAsyncClient, error) {
	// load configuration
	endpoint := taskCtx.GetConfig("JENKINS_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("JENKINS_ENDPOINT is required")
	}
	userRateLimit, err := utils.StrToIntOr(taskCtx.GetConfig("JENKINS_API_REQUESTS_PER_HOUR"), 0)
	if err != nil {
		return nil, err
	}
	username := taskCtx.GetConfig("JENKINS_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("JENKINS_USERNAME is required")
	}
	password := taskCtx.GetConfig("JENKINS_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("JENKINS_PASSWORD is required")
	}
	encodedToken := utils.GetEncodedToken(username, password)
	proxy := taskCtx.GetConfig("JENKINS_PROXY")

	// create synchronize api client so we can calculate api rate limit dynamically
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %v", encodedToken),
	}
	apiClient, err := helper.NewApiClient(endpoint, headers, 0, proxy, taskCtx.GetContext())
	if err != nil {
		return nil, err
	}
	apiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your Username/Password")
		}
		return nil
	})
	// create rate limit calculator
	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: userRateLimit,
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
