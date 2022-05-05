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
	Endpoint string `mapstructure:"GITLAB_ENDPOINT" validate:"required"`
	Auth     string `mapstructure:"GITLAB_AUTH" validate:"required"`
	Proxy    string `mapstructure:"GITLAB_PROXY"`
}

// This object conforms to what the frontend currently expects.
type GitlabResponse struct {
	Name string
	ID   int
	GitlabConnection
}

/*
PUT /plugins/gitlab/connections/:connectionId
*/
func PutConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	gitlabConnection := GitlabConnection{}
	err := mapstructure.Decode(input.Body, &gitlabConnection)
	if err != nil {
		return nil, err
	}
	err = config.SetStruct(gitlabConnection, "mapstructure")
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: "Success"}, nil
}

/*
GET /plugins/gitlab/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-connection is developed.
	gitlabResponse, err := GetConnectionFromEnv()
	response := []GitlabResponse{*gitlabResponse}
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: response}, nil
}

/*
GET /plugins/gitlab/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-connection is developed.)
	gitlabResponse, err := GetConnectionFromEnv()
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: gitlabResponse}, nil
}

func GetConnectionFromEnv() (*GitlabResponse, error) {
	V := config.GetConfig()
	var gitlabConnection GitlabConnection
	err := V.Unmarshal(&gitlabConnection)
	if err != nil {
		return nil, err
	}
	return &GitlabResponse{
		Name:             "Gitlab",
		ID:               1,
		GitlabConnection: gitlabConnection,
	}, nil
}
