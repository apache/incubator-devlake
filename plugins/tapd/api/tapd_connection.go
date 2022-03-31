package api

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/lake/plugins/core"
)

type TapdTestResponse struct {
	Status int `json:"status"`
	Data   struct {
		APIUser     string `json:"api_user"`
		APIPassword string `json:"api_password"`
		RequestIP   string `json:"request_ip"`
	} `json:"data"`
	Info string `json:"info"`
}

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required,url"`
	Auth     string `json:"auth" validate:"required"`
	Proxy    string `json:"proxy"`
}

var vld = validator.New()

/*
POST /plugins/tapd/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// process input
	var params TestConnectionRequest
	err := mapstructure.Decode(input.Body, &params)
	if err != nil {
		return nil, err
	}
	err = vld.Struct(params)
	if err != nil {
		return nil, err
	}
	// verify multiple token in parallel
	// PLEASE NOTE: This works because GitHub API Client rotates tokens on each request
	token := params.Auth
	apiClient, err := helper.NewApiClient(
		params.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", token),
		},
		3*time.Second,
		params.Proxy,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("verify token failed for %s %w", token, err)
	}
	res, err := apiClient.Get("/quickstart/testauth", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("verify token failed for %s %w", token, err)
	}
	resBody := &TapdTestResponse{}
	err = helper.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, err
	}
	// output
	return nil, nil
}
