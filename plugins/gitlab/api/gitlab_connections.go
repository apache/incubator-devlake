package api

import (
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
GET /plugins/gitlab/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	gitlabApiClient := tasks.CreateApiClient()
	headers := gitlabApiClient.GetHeaders()
	endpoint := gitlabApiClient.GetEndpoint()
	if headers["Authorization"] == "Bearer " || endpoint == "" {
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.UnsetConnectionError}}, nil
	}
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
