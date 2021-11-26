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
	githubApiClient := tasks.CreateApiClient(endpoint, tokens)

	// PLEASE NOTE: This works because GitHub API Client rotates tokens on each request
	for i := 0; i < len(tokens); i++ {
		res, err := githubApiClient.Get("user/public_emails", nil, nil)
		if err != nil || res.StatusCode != 200 {
			logger.Error("Error: ", err)
			return &core.ApiResourceOutput{Body: core.TestResult{
				Success: false,
				Message: core.InvalidConnectionError + fmt.Sprintf("Check token: %v or endpoint %v", tokens[i], endpoint),
			}}, nil
		}
		githubApiResponse := &ApiUserPublicEmailResponse{}
		err = core.UnmarshalResponse(res, githubApiResponse)
		if err != nil {
			logger.Error("Error: ", err)
			return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.UnmarshallingError}}, nil
		}
	}
	return &core.ApiResourceOutput{Body: core.TestResult{Success: true, Message: ""}}, nil
}
