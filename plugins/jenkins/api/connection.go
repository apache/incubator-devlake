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
	Endpoint string `mapstructure:"endpoint" validate:"required" env:"JENKINS_ENDPOINT" json:"endpoint"`
	Username string `mapstructure:"username" validate:"required" env:"JENKINS_USERNAME" json:"username"`
	Password string `mapstructure:"password" validate:"required" env:"JENKINS_PASSWORD" json:"password"`
	Proxy    string `mapstructure:"proxy" env:"JENKINS_PROXY" json:"proxy"`
}

type JenkinsResponse struct {
	ID   int
	Name string
	JenkinsConnection
}

/*
PATCH /plugins/jenkins/connections/:connectionId
{
	"Endpoint": "your-endpoint",
	"Username": "your-username",
	"Password": "your-password",
}
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &JenkinsConnection{}, "env")
	if err != nil {
		return nil, err
	}
	// update from request and save to .env
	err = helper.SaveToConifgWithMap(v, connection, input.Body, "env")
	if err != nil {
		return nil, err
	}

	response := JenkinsResponse{
		JenkinsConnection: *connection.(*JenkinsConnection),
		Name:              "Jenkins",
		ID:                1,
	}
	return &core.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

/*
GET /plugins/jenkins/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-source is developed.
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &JenkinsConnection{}, "env")
	if err != nil {
		return nil, err
	}
	response := JenkinsResponse{
		Name:              "Jenkins",
		ID:                1,
		JenkinsConnection: *connection.(*JenkinsConnection),
	}
	return &core.ApiResourceOutput{Body: []JenkinsResponse{response}}, nil
}

/*
GET /plugins/jenkins/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-source is developed.)
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &JenkinsConnection{}, "env")
	if err != nil {
		return nil, err
	}
	response := &JenkinsResponse{
		Name:              "Jenkins",
		ID:                1,
		JenkinsConnection: *connection.(*JenkinsConnection),
	}
	return &core.ApiResourceOutput{Body: response}, nil
}
