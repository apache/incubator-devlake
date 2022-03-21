package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/utils"
)

type JenkinsApiClient struct {
	*helper.ApiAsyncClient
}

func NewJenkinsApiClient(taskCtx core.TaskContext) (*JenkinsApiClient, error) {
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

	// create rate limit calculator
	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: userRateLimit,
	}
	proxy := taskCtx.GetConfig("JENKINS_PROXY")
	asyncApiClient, err := helper.CreateAsyncApiClient(
		taskCtx,
		proxy,
		endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", encodedToken),
		},
		rateLimiter,
	)
	if err != nil {
		return nil, err
	}

	jenkinsApiClient := &JenkinsApiClient{
		asyncApiClient,
	}

	jenkinsApiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your Username/Password")
		}
		return nil
	})
	return jenkinsApiClient, nil
}
