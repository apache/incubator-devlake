package api

import (
	"fmt"
	"strings"

	"github.com/merico-dev/lake/logger"
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
GET /plugins/github/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	ValidationResult := core.ValidateParams(input, []string{"endpoint", "auth"})
	if !ValidationResult.Success {
		return &core.ApiResourceOutput{Body: ValidationResult}, nil
	}
	endpoint := input.Query.Get("endpoint")
	auth := input.Query.Get("auth")
	tokens := strings.Split(auth, ",")

	// PLEASE NOTE: This works because GitHub API Client rotates tokens on each request
	results := make(chan error)
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		i := i
		go func() {
			githubApiClient := tasks.CreateApiClient(endpoint, []string{token})
			res, err := githubApiClient.Get("user/public_emails", nil, nil)
			if err != nil || res.StatusCode != 200 {
				logger.Error("Error: ", err)
				results <- fmt.Errorf("invalid token #%v %s", i, token)
				return
			}
			githubApiResponse := &ApiUserPublicEmailResponse{}
			err = core.UnmarshalResponse(res, githubApiResponse)
			if err != nil {
				logger.Error("Error: ", err)
				results <- fmt.Errorf("invalid token #%v %s", i, token)
			} else {
				results <- nil
			}
		}()
	}

	println("length of tokens", len(tokens))
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
	return &core.ApiResourceOutput{Body: core.TestResult{Success: len(msgs) == 0, Message: strings.Join(msgs, "\n")}}, nil
}
