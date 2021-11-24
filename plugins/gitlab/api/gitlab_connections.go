package api

import (
	"errors"
	"fmt"

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

	res, err := gitlabApiClient.Get("user", nil, nil)
	if err != nil {
		logger.Error("Error: ", err)
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("Your token is invalid for this URL.")
	}
	fmt.Println("KEVIN >>> res.Status", res.Status)
	gitlabApiResponse := &ApiUserResponse{}
	err = core.UnmarshalResponse(res, gitlabApiResponse)
	if err != nil {
		logger.Error("Error: ", err)
		return nil, err
	}
	fmt.Println("KEVIN >>> gitlabApiResponse.Name", gitlabApiResponse.Name)
	return &core.ApiResourceOutput{Body: "Token is valid for the URL!"}, nil
}
