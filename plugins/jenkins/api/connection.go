package api

import (
	"fmt"
	"github.com/merico-dev/lake/config"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/utils"
	"github.com/mitchellh/mapstructure"
)

var vld = validator.New()

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Proxy    string `json:"proxy"`
}

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// decode
	var err error
	var connection TestConnectionRequest
	err = mapstructure.Decode(input.Body, &connection)
	if err != nil {
		return nil, err
	}
	// validate
	err = vld.Struct(connection)
	if err != nil {
		return nil, err
	}
	// test connection
	encodedToken := utils.GetEncodedToken(connection.Username, connection.Password)
	apiClient, err := helper.NewApiClient(
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", encodedToken),
		},
		3*time.Second,
		connection.Proxy,
		nil,
	)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("", nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil, nil
}

// This object conforms to what the frontend currently sends.
type JenkinsConnection struct {
	Endpoint string `mapstructure:"JENKINS_ENDPOINT" validate:"required"`
	Username string `mapstructure:"JENKINS_USERNAME" validate:"required"`
	Password string `mapstructure:"JENKINS_PASSWORD" validate:"required"`
	Proxy    string `mapstructure:"JENKINS_PROXY"`
}

type JenkinsResponse struct {
	ID   int
	Name string
	JenkinsConnection
}

/*
POST /plugins/jenkins/connections
{
	"Endpoint": "your-endpoint",
	"Username": "your-username",
	"Password": "your-password",
}
*/
func PostConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// TODO: For now, Jenkins does not support sources but it will in the future.

	return PutConnection(input)
}

/*
PUT /plugins/jenkins/connections/:connectionId
{
	"Endpoint": "your-endpoint",
	"Username": "your-username",
	"Password": "your-password",
}
*/
func PutConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jenkinsConnection := JenkinsConnection{}
	err := mapstructure.Decode(input.Body, &jenkinsConnection)
	if err != nil {
		return nil, err
	}
	err = config.SetStruct(jenkinsConnection, "mapstructure")
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: "Success"}, nil
}

/*
GET /plugins/jenkins/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-source is developed.
	jenkinsResponse, err := GetConnectionFromEnv()
	response := []JenkinsResponse{*jenkinsResponse}
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: response}, nil
}

/*
GET /plugins/jenkins/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-source is developed.)
	jenkinsResponse, err := GetConnectionFromEnv()
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jenkinsResponse}, nil
}

func GetConnectionFromEnv() (*JenkinsResponse, error) {
	v := config.GetConfig()
	var jenkinsConnection JenkinsConnection
	err := v.Unmarshal(&jenkinsConnection)
	if err != nil {
		return nil, err
	}
	return &JenkinsResponse{
		Name:              "Jenkins",
		ID:                1,
		JenkinsConnection: jenkinsConnection,
	}, nil
}
