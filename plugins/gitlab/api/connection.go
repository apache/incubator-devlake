package api

import (
	"fmt"
	"github.com/merico-dev/lake/config"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/mitchellh/mapstructure"
)

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string
}

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required"`
	Auth     string `json:"auth" validate:"required"`
	Proxy    string `json:"proxy"`
}

var vld = validator.New()

/*
POST /plugins/gitlab/test
*/
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
	apiClient, err := helper.NewApiClient(
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", connection.Auth),
		},
		3*time.Second,
		connection.Proxy,
		nil,
	)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("user", nil, nil)
	if err != nil {
		return nil, err
	}
	resBody := &ApiUserResponse{}
	err = helper.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil, nil
}

// This object conforms to what the frontend currently sends.
type GitlabConnection struct {
	Endpoint string `mapstructure:"endpoint" validate:"required" env:"GITLAB_ENDPOINT" json:"endpoint"`
	Auth     string `mapstructure:"auth" validate:"required" env:"GITLAB_AUTH"  json:"auth"`
	Proxy    string `mapstructure:"proxy" env:"GITLAB_PROXY" json:"proxy"`
}

// This object conforms to what the frontend currently expects.
type GitlabResponse struct {
	Name string
	ID   int
	GitlabConnection
}

/*
PATCH /plugins/gitlab/connections/:connectionId
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &GitlabConnection{}, "env")
	if err != nil {
		return nil, err
	}
	// update from request and save to .env
	err = helper.SaveToConifgWithMap(v, connection, input.Body, "env")
	if err != nil {
		return nil, err
	}

	response := GitlabResponse{
		GitlabConnection: *connection.(*GitlabConnection),
		Name:             "Gitlab",
		ID:               1,
	}
	return &core.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

/*
GET /plugins/gitlab/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-connection is developed.
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &GitlabConnection{}, "env")
	if err != nil {
		return nil, err
	}
	response := GitlabResponse{
		GitlabConnection: *connection.(*GitlabConnection),
		Name:             "Gitlab",
		ID:               1,
	}

	return &core.ApiResourceOutput{Body: []GitlabResponse{response}}, nil
}

/*
GET /plugins/gitlab/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-connection is developed.)
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &GitlabConnection{}, "env")
	if err != nil {
		return nil, err
	}

	response := &GitlabResponse{
		GitlabConnection: *connection.(*GitlabConnection),
		Name:             "Gitlab",
		ID:               1,
	}
	return &core.ApiResourceOutput{Body: response}, nil
}
