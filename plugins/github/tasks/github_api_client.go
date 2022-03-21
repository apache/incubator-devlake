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

type GithubApiClient struct {
	*helper.ApiAsyncClient
	// This is for multiple token functionality so we can loop through an array of tokens.
	tokens     []string
	tokenIndex int
}

func NewGithubApiClient(taskCtx core.TaskContext) (*GithubApiClient, error) {
	// load configuration
	endpoint := taskCtx.GetConfig("GITHUB_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("endpint is required")
	}
	userRateLimit, err := utils.StrToIntOr(taskCtx.GetConfig("GITHUB_API_REQUESTS_PER_HOUR"), 0)
	if err != nil {
		return nil, err
	}
	auth := taskCtx.GetConfig("GITHUB_AUTH")
	if auth == "" {
		return nil, fmt.Errorf("GITHUB_AUTH is required")
	}
	tokens := strings.Split(auth, ",")

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
	proxy := taskCtx.GetConfig("GITHUB_PROXY")
	asyncApiClient, err := helper.CreateAsyncApiClient(
		taskCtx,
		proxy,
		endpoint,
		nil,
		rateLimiter,
	)
	if err != nil {
		return nil, err
	}

	githubApiClient := &GithubApiClient{
		asyncApiClient,
		tokens,
		0,
	}
	// Rotates token on each request.
	githubApiClient.SetBeforeFunction(func(req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", githubApiClient.tokens[githubApiClient.tokenIndex]))
		// Set next token index
		githubApiClient.tokenIndex = (githubApiClient.tokenIndex + 1) % len(githubApiClient.tokens)
		return nil
	})
	githubApiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your Token configuration")
		}
		return nil
	})

	return githubApiClient, nil
}
