package api

import (
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/tasks"
)

type ApiMyselfResponse struct {
	AccountId   string
	DisplayName string
}

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	ValidationResult := core.ValidateParams(input, []string{"endpoint", "auth"})
	if !ValidationResult.Success {
		return &core.ApiResourceOutput{Body: ValidationResult}, nil
	}
	endpoint := input.Query.Get("endpoint")
	auth := input.Query.Get("auth")
	jiraApiClient := tasks.NewJiraApiClient(endpoint, auth)

	res, err := jiraApiClient.Get("api/3/myself", nil, nil)
	if err != nil || res.StatusCode != 200 {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{
			Success: false,
			Message: core.InvalidConnectionError}}, nil
	}

	myselfFromApi := &ApiMyselfResponse{}

	err = core.UnmarshalResponse(res, myselfFromApi)
	if err != nil {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.UnmarshallingError}}, nil
	}
	return &core.ApiResourceOutput{Body: core.TestResult{Success: true, Message: ""}}, nil
}
