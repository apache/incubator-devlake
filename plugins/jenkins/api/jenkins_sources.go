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
	Username string
	Password string
	Endpoint string
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
	V := config.LoadConfigFile()
	V.Set("JENKINS_ENDPOINT", jenkinsSource.Endpoint)
	V.Set("JENKINS_USERNAME", jenkinsSource.Username)
	V.Set("JENKINS_PASSWORD", jenkinsSource.Password)
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
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jenkinsSources}, nil
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

func GetSourceFromEnv() (*[1]JenkinsSource, error) {
	V := config.LoadConfigFile()
	var configJson JenkinsConfig
	err := V.Unmarshal(&configJson)
	if err != nil {
		return nil, err
	}
	return &[1]JenkinsSource{{
		Endpoint: configJson.JENKINS_ENDPOINT,
		Username: configJson.JENKINS_USERNAME,
		Password: configJson.JENKINS_PASSWORD,
	}}, nil
}
