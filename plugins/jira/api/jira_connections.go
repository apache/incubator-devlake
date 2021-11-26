package api

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/tasks"
)

type ApiMyselfResponse struct {
	AccountId   string
	DisplayName string
}

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraSource, err := findSourceByInputQuery(input)
	if err != nil {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.SourceIdError}}, nil
	}
	jiraApiClient := tasks.NewJiraApiClient(jiraSource.Endpoint, jiraSource.BasicAuthEncoded)

	res, err := jiraApiClient.Get("api/3/myself", nil, nil)
	if err != nil || res.StatusCode != 200 {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{
			Success: false, 
			Message: fmt.Sprintf("Your connection configuration is invalid for this source: %v", jiraSource.Name)}}, nil
	}

	myselfFromApi := &ApiMyselfResponse{}

	err = core.UnmarshalResponse(res, myselfFromApi)
	if err != nil {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.UnmarshallingError}}, nil
	}
	return &core.ApiResourceOutput{Body: core.TestResult{Success: true, Message: ""}}, nil
}
