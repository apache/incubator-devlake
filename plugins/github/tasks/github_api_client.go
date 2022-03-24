package tasks

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/utils"
)

func CreateApiClient(taskCtx core.TaskContext) (*helper.ApiAsyncClient, error) {
	// load configuration
	endpoint := taskCtx.GetConfig("GITHUB_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("endpint is required")
	}
	proxy := taskCtx.GetConfig("GITHUB_PROXY")
	userRateLimit, err := utils.StrToIntOr(taskCtx.GetConfig("GITHUB_API_REQUESTS_PER_HOUR"), 0)
	if err != nil {
		return nil, err
	}
	auth := taskCtx.GetConfig("GITHUB_AUTH")
	if auth == "" {
		return nil, fmt.Errorf("GITHUB_AUTH is required")
	}
	tokens := strings.Split(auth, ",")
	tokenIndex := 0
	// create synchronize api client so we can calculate api rate limit dynamically
	apiClient, err := helper.NewApiClient(endpoint, nil, 0, proxy, taskCtx.GetContext())
	if err != nil {
		return nil, err
	}
	// Rotates token on each request.
	apiClient.SetBeforeFunction(func(req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokens[tokenIndex]))
		// Set next token index
		tokenIndex = (tokenIndex + 1) % len(tokens)
		return nil
	})
	apiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your Token configuration")
		}
		return nil
	})

	// create rate limit calculator
	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: userRateLimit,
		DynamicRateLimitPerHour: func(res *http.Response) (int, time.Duration, error) {
			/* calculate by number of remaining requests
			remaining, err := strconv.Atoi(res.Header.Get("X-RateLimit-Remaining"))
			if err != nil {
				return 0,0, fmt.Errorf("failed to parse X-RateLimit-Remaining header: %w", err)
			}
			reset, err := strconv.Atoi(res.Header.Get("X-RateLimit-Reset"))
			if err != nil {
				return 0, 0, fmt.Errorf("failed to parse X-RateLimit-Reset header: %w", err)
			}
			date, err := http.ParseTime(res.Header.Get("Date"))
			if err != nil {
				return 0, 0, fmt.Errorf("failed to parse Date header: %w", err)
			}
			return remaining * len(tokens), time.Unix(int64(reset), 0).Sub(date), nil
			*/
			rateLimit, err := strconv.Atoi(res.Header.Get("X-RateLimit-Limit"))
			if err != nil {
				return 0, 0, fmt.Errorf("failed to parse X-RateLimit-Limit header: %w", err)
			}
			// even though different token could have different rate limit, but it is hard to support it
			// so, we calculate the rate limit of a single token, and presume all tokens are the same, to
			// simplify the algorithm for now
			// TODO: consider different token has different rate-limit
			return rateLimit * len(tokens), 1 * time.Hour, nil

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
