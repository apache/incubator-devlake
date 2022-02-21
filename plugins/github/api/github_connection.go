package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/tasks"
)

type ApiUserPublicEmailResponse []PublicEmail

// Using Public Email because it requires authentication, and it is public information anyway.
// We're not using email information for anything here.
type PublicEmail struct {
	Email      string
	Primary    bool
	Verified   bool
	Visibility string
}

/*
POST /plugins/github/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// process input
	ValidationResult := core.ValidateParams(input, []string{"endpoint", "auth"})
	if !ValidationResult.Success {
		return &core.ApiResourceOutput{Body: ValidationResult}, nil
	}
	endpoint := input.Body["endpoint"].(string)
	auth := input.Body["auth"].(string)
	proxy := input.Body["proxy"].(string)
	tokens := strings.Split(auth, ",")

	// verify multiple token in parallel
	// PLEASE NOTE: This works because GitHub API Client rotates tokens on each request
	results := make(chan error)
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		i := i
		go func() {
			githubApiClient := tasks.CreateApiClient(endpoint, []string{token}, nil)
			githubApiClient.SetTimeout(3 * time.Second)
			if proxy != "" {
				githubApiClient.SetProxy(proxy)
			}
			res, err := githubApiClient.Get("user/public_emails", nil, nil)
			if err != nil {
				results <- fmt.Errorf("verify token failed for #%v %s %w", i, token, err)
				return
			}
			githubApiResponse := &ApiUserPublicEmailResponse{}
			err = core.UnmarshalResponse(res, githubApiResponse)
			if err != nil {
				results <- fmt.Errorf("verify token failed for #%v %s %w", i, token, err)
			} else {
				results <- nil
			}
		}()
	}

	// collect verification results
	msgs := make([]string, 0)
	i := 0
	for err := range results {
		if err != nil {
			msgs = append(msgs, err.Error())
		}
		i++
		if i == len(tokens) {
			close(results)
		}
	}

	// output
	return &core.ApiResourceOutput{
		Body: core.TestResult{
			Success: len(msgs) == 0,
			Message: strings.Join(msgs, "\n"),
		},
	}, nil
}
