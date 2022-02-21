package api

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/tasks"
)

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string
}

/*
POST /plugins/gitlab/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	gitlabApiClient := tasks.CreateApiClient()

	ValidationResult := core.ValidateParams(input, []string{"endpoint", "auth"})
	if !ValidationResult.Success {
		return &core.ApiResourceOutput{Body: ValidationResult}, nil
	}
	endpoint := input.Body["endpoint"].(string)
	proxy := input.Body["proxy"].(string)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", input.Body["auth"].(string)),
	}
	gitlabApiClient.SetEndpoint(endpoint)
	gitlabApiClient.SetHeaders(headers)
	gitlabApiClient.SetProxy(proxy)
	gitlabApiClient.SetTimeout(3 * time.Second)

	res, err := gitlabApiClient.Get("user", nil, nil)
	if err != nil || res.StatusCode != 200 {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.InvalidConnectionError}}, nil
	}

	gitlabApiResponse := &ApiUserResponse{}

	err = core.UnmarshalResponse(res, gitlabApiResponse)
	if err != nil {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.UnmarshallingError}}, nil
	}
	return &core.ApiResourceOutput{Body: core.TestResult{Success: true, Message: ""}}, nil
}
