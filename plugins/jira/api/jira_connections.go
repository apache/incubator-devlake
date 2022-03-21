package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/mitchellh/mapstructure"
)

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint"`
	Auth     string `json:"auth"`
	Proxy    string `json:"proxy"`
}

var vld = validator.New()

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
	res, err := apiClient.Get("api/2/serverInfo", nil, nil)
	if err != nil {
		return nil, err
	}
	resBody := &models.JiraServerInfo{}
	err = helper.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, err
	}
	// check version
	if resBody.DeploymentType == models.DeploymentServer {
		// only support 8.x.x or higher
		if versions := resBody.VersionNumbers; len(versions) == 3 && versions[0] < 8 {
			return nil, fmt.Errorf("Support JIRA Server 8+ only")
		}
	}
	return &core.ApiResourceOutput{
		Status: res.StatusCode,
		Body: &core.TestResult{
			Success: res.StatusCode == http.StatusOK,
			Message: "success",
		},
	}, nil
}
