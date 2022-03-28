package tasks

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/utils"
)

func NewGitlabApiClient(taskCtx core.TaskContext) (*helper.ApiAsyncClient, error) {
	// load configuration
	endpoint := taskCtx.GetConfig("GITLAB_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("endpint is required")
	}
	userRateLimit, err := utils.StrToIntOr(taskCtx.GetConfig("GITLAB_API_REQUESTS_PER_HOUR"), 0)
	if err != nil {
		return nil, err
	}
	auth := taskCtx.GetConfig("GITLAB_AUTH")
	if auth == "" {
		return nil, fmt.Errorf("GITLAB_AUTH is required")
	}
	proxy := taskCtx.GetConfig("GITLAB_PROXY")

	// create synchronize api client so we can calculate api rate limit dynamically
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", auth),
	}
	apiClient, err := helper.NewApiClient(endpoint, headers, 0, proxy, taskCtx.GetContext())
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
		UserRateLimitPerHour: userRateLimit,
		DynamicRateLimit: func(res *http.Response) (int, time.Duration, error) {
			rateLimitHeader := res.Header.Get("RateLimit-Limit")
			if rateLimitHeader == "" {
				// unlimited
				return 0, 0, nil
			}
			rateLimit, err := strconv.Atoi(rateLimitHeader)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to parse RateLimit-Limit header: %w", err)
			}
			// seems like gitlab rate limit is on minute basis
			return rateLimit, 1 * time.Minute, nil
		},
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
