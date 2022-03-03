package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
)

// This object conforms to what the frontend currently sends.
type GitlabSource struct {
	Endpoint string `mapstructure:"GITLAB_ENDPOINT" validate:"required"`
	Auth     string `mapstructure:"GITLAB_AUTH" validate:"required"`
	Proxy    string `mapstructure:"GITLAB_PROXY"`
}

// This object conforms to what the frontend currently expects.
type GitlabResponse struct {
	Name string
	ID   int
	GitlabSource
}

/*
PUT /plugins/gitlab/sources/:sourceId
*/
func PutSource(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	gitlabSource := GitlabSource{}
	err := mapstructure.Decode(input.Body, &gitlabSource)
	if err != nil {
		return nil, err
	}
	err = config.SetStruct(gitlabSource, "mapstructure")
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: "Success"}, nil
}

/*
GET /plugins/gitlab/sources
*/
func ListSources(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-source is developed.
	gitlabResponse, err := GetSourceFromEnv()
	response := []GitlabResponse{*gitlabResponse}
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: response}, nil
}

/*
GET /plugins/gitlab/sources/:sourceId
*/
func GetSource(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-source is developed.)
	gitlabResponse, err := GetSourceFromEnv()
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: gitlabResponse}, nil
}

func GetSourceFromEnv() (*GitlabResponse, error) {
	V := config.GetConfig()
	var gitlabSource GitlabSource
	err := V.Unmarshal(&gitlabSource)
	if err != nil {
		return nil, err
	}
	return &GitlabResponse{
		Name:         "Gitlab",
		ID:           1,
		GitlabSource: gitlabSource,
	}, nil
}
