package api

import (
	"fmt"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/jenkins/models"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/utils"
	"github.com/mitchellh/mapstructure"
)

var vld = validator.New()

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// decode
	var err error
	var connection models.TestConnectionRequest
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
	connection := &models.JenkinsConnection{}
	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	// update from request and save to .env
	err = helper.DecodeStruct(v, connection, input.Body, "env")
	if err != nil {
		return nil, err
	}
	err = v.WriteConfig()
	if err != nil {
		return nil, err
	}
	response := models.JenkinsResponse{
		JenkinsConnection: *connection,
		Name:              "Jenkins",
		ID:                1,
	}
	return &core.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

/*
GET /plugins/jenkins/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-connection is developed.
	v := config.GetConfig()
	connection := &models.JenkinsConnection{}

	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	response := models.JenkinsResponse{
		JenkinsConnection: *connection,
		Name:              "Jenkins",
		ID:                1,
	}

	return &core.ApiResourceOutput{Body: []models.JenkinsResponse{response}}, nil
}

/*
GET /plugins/jenkins/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-connection is developed.)
	v := config.GetConfig()
	connection := &models.JenkinsConnection{}
	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	response := &models.JenkinsResponse{
		JenkinsConnection: *connection,
		Name:              "Jenkins",
		ID:                1,
	}
	return &core.ApiResourceOutput{Body: response}, nil
}
