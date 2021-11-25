package api

import (
	"fmt"
	"strings"

	"github.com/merico-dev/lake/config"
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
	endpoint := config.V.GetString("GITHUB_ENDPOINT")
	configTokensString := config.V.GetString("GITHUB_AUTH")
	tokens := strings.Split(configTokensString, ",")
	githubApiClient := tasks.CreateApiClient(endpoint, tokens)

	// PLEASE NOTE: This works because GitHub API Client rotates tokens on each request
	for i := 0; i < len(tokens); i++ {
		res, err := githubApiClient.Get("/user/public_emails", nil, nil)
		if err != nil {
			logger.Error("Error: ", err)
			return nil, err
		}
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("Invalid token: %v. Please ensure your tokens are correct.", tokens[i])
		}
		githubApiResponse := &ApiUserPublicEmailResponse{}
		err = core.UnmarshalResponse(res, githubApiResponse)
		if err != nil {
			logger.Error("Error: ", err)
			return nil, err
		}
	}
	return &core.ApiResourceOutput{Body: true}, nil
}
