package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
)

// This object conforms to what the frontend currently sends.
type GithubSource struct {
	Endpoint string `mapstructure:"GITHUB_ENDPOINT" validate:"required"`
	Auth     string `mapstructure:"GITHUB_AUTH" validate:"required"`
	Proxy    string `mapstructure:"GITHUB_PROXY"`

	PrType               string `mapstructure:"GITHUB_PR_TYPE"`
	PrComponent          string `mapstructure:"GITHUB_PR_COMPONENT"`
	IssueSeverity        string `mapstructure:"GITHUB_ISSUE_SEVERITY"`
	IssuePriority        string `mapstructure:"GITHUB_ISSUE_PRIORITY"`
	IssueComponent       string `mapstructure:"GITHUB_ISSUE_COMPONENT"`
	IssueTypeBug         string `mapstructure:"GITHUB_ISSUE_TYPE_BUG"`
	IssueTypeIncident    string `mapstructure:"GITHUB_ISSUE_TYPE_INCIDENT"`
	IssueTypeRequirement string `mapstructure:"GITHUB_ISSUE_TYPE_REQUIREMENT"`
}

// This object conforms to what the frontend currently expects.
type GithubResponse struct {
	Name string
	ID   int

	GithubSource
}

/*
PUT /plugins/github/sources/:sourceId
*/
func PutSource(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	githubSource := GithubSource{}
	err := mapstructure.Decode(input.Body, &githubSource)
	if err != nil {
		return nil, err
	}

	err = config.SetStruct(githubSource, "mapstructure")
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: "Success"}, nil
}

/*
GET /plugins/github/sources
*/
func ListSources(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-source is developed.
	githubSources, err := GetSourceFromEnv()
	response := []GithubResponse{*githubSources}
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: response}, nil
}

/*
GET /plugins/github/sources/:sourceId
*/
func GetSource(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-source is developed.)
	githubSources, err := GetSourceFromEnv()
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: githubSources}, nil
}

func GetSourceFromEnv() (*GithubResponse, error) {
	v := config.GetConfig()
	var githubSource GithubSource
	err := v.Unmarshal(&githubSource)
	if err != nil {
		return nil, err
	}

	return &GithubResponse{
		Name:         "Github",
		ID:           1,
		GithubSource: githubSource,
	}, nil
}
