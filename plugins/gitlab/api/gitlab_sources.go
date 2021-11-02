package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var V *viper.Viper

type GitlabConfig struct {
	GITLAB_ENDPOINT string `mapstructure:"GITLAB_ENDPOINT"`
	GITLAB_AUTH     string `mapstructure:"GITLAB_AUTH"`
}
type GitlabSource struct {
	Endpoint string
	Auth     string
}

/*
PUT /plugins/gitlab/sources/:sourceId
{
	"Endpoint": "",
	"Auth": ""

}
*/
func PutSource(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	gitlabSource := GitlabSource{}
	err := mapstructure.Decode(input.Body, &gitlabSource)
	if err != nil {
		return nil, err
	}
	V := config.LoadConfigFile()
	V.Set("GITLAB_ENDPOINT", gitlabSource.Endpoint)
	V.Set("GITLAB_AUTH", gitlabSource.Auth)
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
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: gitlabSources}, nil
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

func GetSourceFromEnv() (*[1]GitlabSource, error) {
	V := config.LoadConfigFile()
	var configJson GitlabConfig
	err := V.Unmarshal(&configJson)
	if err != nil {
		return nil, err
	}
	return &[1]GitlabSource{{
		Endpoint: configJson.GITLAB_ENDPOINT,
		Auth:     configJson.GITLAB_AUTH,
	}}, nil
}
