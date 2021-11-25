package api

import (
	"fmt"

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
		return nil, err
	}
	jiraApiClient := tasks.NewJiraApiClient(jiraSource.Endpoint, jiraSource.BasicAuthEncoded)

	res, err := jiraApiClient.Get("api/3/myself", nil, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Your connection configuration is invalid for this source: %v", jiraSource.Name)
	}
	myselfFromApi := &ApiMyselfResponse{}

	err = core.UnmarshalResponse(res, myselfFromApi)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: true}, nil
}
