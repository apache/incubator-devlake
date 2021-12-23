package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var V *viper.Viper

type JenkinsConfig struct {
	JENKINS_ENDPOINT string `mapstructure:"JENKINS_ENDPOINT"`
	JENKINS_USERNAME string `mapstructure:"JENKINS_USERNAME"`
	JENKINS_PASSWORD string `mapstructure:"JENKINS_PASSWORD"`
}

type JenkinsSource struct {
	ID       int
	Username string
	Password string
	Endpoint string
	Name     string
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
	jenkinsSource := JenkinsConfig{}
	err := mapstructure.Decode(input.Body, &jenkinsSource)
	if err != nil {
		return nil, err
	}
	V, err := config.LoadConfigFile()
	if err != nil {
		return nil, err
	}
	V.Set("JENKINS_ENDPOINT", jenkinsSource.JENKINS_ENDPOINT)
	V.Set("JENKINS_USERNAME", jenkinsSource.JENKINS_USERNAME)
	V.Set("JENKINS_PASSWORD", jenkinsSource.JENKINS_PASSWORD)
	err = V.WriteConfig()
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
	jenkinsSources, err := GetSourceFromEnv()
	response := []JenkinsSource{*jenkinsSources}
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
	jenkinsSources, err := GetSourceFromEnv()
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jenkinsSources}, nil
}

func GetSourceFromEnv() (*JenkinsSource, error) {
	configJson, err := config.GetConfigJson()
	if err != nil {
		return nil, err
	}
	return &JenkinsSource{
		Endpoint: configJson.JENKINS_ENDPOINT,
		Username: configJson.JENKINS_USERNAME,
		Password: configJson.JENKINS_PASSWORD,
		// The UI relies on a source ID here but we will hardcode it until the sources work is done for Jenkins
		ID:   1,
		Name: "Jenkins",
	}, nil
}
