package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/lake/plugins/core"
)

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
		return nil, err
	}
	if res.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("verify token failed for %s", token)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	// output
	return nil, nil
}
