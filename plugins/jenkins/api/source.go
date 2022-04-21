package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
)

// This object conforms to what the frontend currently sends.
type JenkinsSource struct {
	Endpoint string `mapstructure:"JENKINS_ENDPOINT" validate:"required"`
	Username string `mapstructure:"JENKINS_USERNAME" validate:"required"`
	Password string `mapstructure:"JENKINS_PASSWORD" validate:"required"`
	Proxy    string `mapstructure:"JENKINS_PROXY"`
}

type JenkinsResponse struct {
	ID   int
	Name string
	JenkinsSource
}

/*
POST /plugins/jenkins/sources
{
	"Endpoint": "your-endpoint",
	"Username": "your-username",
	"Password": "your-password",
}
*/
func PostSource(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// TODO: For now, Jenkins does not support sources but it will in the future.

	return PutSource(input)
}

/*
PUT /plugins/jenkins/sources/:sourceId
{
	"Endpoint": "your-endpoint",
	"Username": "your-username",
	"Password": "your-password",
}
*/
func PutSource(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jenkinsSource := JenkinsSource{}
	err := mapstructure.Decode(input.Body, &jenkinsSource)
	if err != nil {
		return nil, err
	}
	err = config.SetStruct(jenkinsSource, "mapstructure")
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: "Success"}, nil
}

/*
GET /plugins/jenkins/sources
*/
func ListSources(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-source is developed.
	jenkinsResponse, err := GetSourceFromEnv()
	response := []JenkinsResponse{*jenkinsResponse}
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: response}, nil
}

/*
GET /plugins/jenkins/sources/:sourceId
*/
func GetSource(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-source is developed.)
	jenkinsResponse, err := GetSourceFromEnv()
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jenkinsResponse}, nil
}

func GetSourceFromEnv() (*JenkinsResponse, error) {
	v := config.GetConfig()
	var jenkinsSource JenkinsSource
	err := v.Unmarshal(&jenkinsSource)
	if err != nil {
		return nil, err
	}
	return &JenkinsResponse{
		Name:          "Jenkins",
		ID:            1,
		JenkinsSource: jenkinsSource,
	}, nil
}
