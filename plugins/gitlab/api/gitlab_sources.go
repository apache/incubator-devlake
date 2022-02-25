package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
)

type GitlabConfig struct {
	GITLAB_ENDPOINT            string `mapstructure:"GITLAB_ENDPOINT"`
	GITLAB_AUTH                string `mapstructure:"GITLAB_AUTH"`
	JIRA_BOARD_GITLAB_PROJECTS string `mapstructure:"JIRA_BOARD_GITLAB_PROJECTS"`
}

// This object conforms to what the frontend currently sends.
type GitlabSource struct {
	GITLAB_ENDPOINT            string
	GITLAB_AUTH                string
	JIRA_BOARD_GITLAB_PROJECTS string
}

// This object conforms to what the frontend currently expects.
type GitlabResponse struct {
	Endpoint                   string
	Auth                       string
	Name                       string
	ID                         int
	JIRA_BOARD_GITLAB_PROJECTS string
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
	V := config.GetConfig()
	if gitlabSource.GITLAB_ENDPOINT != "" {
		V.Set("GITLAB_ENDPOINT", gitlabSource.GITLAB_ENDPOINT)
	}
	if gitlabSource.GITLAB_AUTH != "" {
		V.Set("GITLAB_AUTH", gitlabSource.GITLAB_AUTH)
	}
	if gitlabSource.JIRA_BOARD_GITLAB_PROJECTS != "" {
		V.Set("JIRA_BOARD_GITLAB_PROJECTS", gitlabSource.JIRA_BOARD_GITLAB_PROJECTS)
	}
	err = V.WriteConfig()
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
	gitlabSources, err := GetSourceFromEnv()
	response := []GitlabResponse{*gitlabSources}
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
	gitlabSources, err := GetSourceFromEnv()
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: gitlabSources}, nil
}

func GetSourceFromEnv() (*GitlabResponse, error) {
	V := config.GetConfig()
	var configJson GitlabConfig
	err := V.Unmarshal(&configJson)
	if err != nil {
		return nil, err
	}
	return &GitlabResponse{
		Endpoint:                   configJson.GITLAB_ENDPOINT,
		Auth:                       configJson.GITLAB_AUTH,
		Name:                       "Gitlab",
		ID:                         1,
		JIRA_BOARD_GITLAB_PROJECTS: configJson.JIRA_BOARD_GITLAB_PROJECTS,
	}, nil
}
