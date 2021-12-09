package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
)

type GithubConfig struct {
	GITHUB_ENDPOINT string `mapstructure:"GITHUB_ENDPOINT"`
	GITHUB_AUTH     string `mapstructure:"GITHUB_AUTH"`
	GITHUB_PROXY    string `mapstructure:"GITHUB_PROXY"`
}

// This object conforms to what the frontend currently sends.
type GithubSource struct {
	GITHUB_ENDPOINT string
	GITHUB_AUTH     string
	GITHUB_PROXY    string
}

// This object conforms to what the frontend currently expects.
type GithubResponse struct {
	Endpoint string
	Auth     string
	Name     string
	ID       int
	Proxy    string `json:"proxy"`
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
	V := config.LoadConfigFile()
	if githubSource.GITHUB_ENDPOINT != "" {
		V.Set("GITHUB_ENDPOINT", githubSource.GITHUB_ENDPOINT)
	}
	if githubSource.GITHUB_AUTH != "" {
		V.Set("GITHUB_AUTH", githubSource.GITHUB_AUTH)
	}
	V.Set("GITHUB_PROXY", githubSource.GITHUB_PROXY)
	err = V.WriteConfig()
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
	V := config.LoadConfigFile()
	var configJson GithubConfig
	err := V.Unmarshal(&configJson)
	if err != nil {
		return nil, err
	}
	return &GithubResponse{
		Endpoint: configJson.GITHUB_ENDPOINT,
		Auth:     configJson.GITHUB_AUTH,
		Name:     "Github",
		ID:       1,
		Proxy:    configJson.GITHUB_PROXY,
	}, nil
}
