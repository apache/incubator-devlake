package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/utils"
)

type GithubApiClient struct {
	tokenIndex int
	// This is for multiple token functionality so we can loop through an array of tokens.
	tokens []string
	core.ApiClient
}

func NewGithubApiClient(endpoint string, tokens []string, proxy string, ctx context.Context, scheduler *utils.WorkerScheduler, logger core.Logger) *GithubApiClient {
	githubApiClient := &GithubApiClient{}
	githubApiClient.tokenIndex = 0
	githubApiClient.tokens = tokens
	// Rotates token on each request.
	githubApiClient.SetBeforeFunction(func(req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", githubApiClient.tokens[githubApiClient.tokenIndex]))
		// Set next token index
		githubApiClient.tokenIndex = (githubApiClient.tokenIndex + 1) % len(githubApiClient.tokens)
		return nil
	})
	githubApiClient.Setup(
		endpoint,
		map[string]string{},
		50*time.Second,
		3,
		scheduler,
	)
	if ctx != nil {
		githubApiClient.SetContext(ctx)
	}
	if proxy != "" {
		err := githubApiClient.SetProxy(proxy)
		if err != nil {
			panic(err)
		}
	}

	githubApiClient.SetLogger(logger)
	return githubApiClient
}
