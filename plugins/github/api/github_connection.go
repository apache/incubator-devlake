package api

import (
	"strings"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/tasks"
)

type ApiMeResponse struct {
	Name     string `json:"name"`
	GithubId int    `json:"id"`
	HTMLUrl  string `json:"html_url"`
}

/*
GET /plugins/github/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	endpoint := config.V.GetString("GITHUB_ENDPOINT")
	configTokensString := config.V.GetString("GITHUB_AUTH")
	tokens := strings.Split(configTokensString, ",")
	githubApiClient := tasks.CreateApiClient(endpoint, tokens)

	res, err := githubApiClient.Get("/users/me", nil, nil)
	if err != nil {
		logger.Error("Error: ", err)
		return nil, err
	}

	githubApiResponse := &ApiMeResponse{}
	err = core.UnmarshalResponse(res, githubApiResponse)
	if err != nil {
		logger.Error("Error: ", err)
		return nil, err
	}

	return &core.ApiResourceOutput{Body: githubApiResponse}, nil
}
