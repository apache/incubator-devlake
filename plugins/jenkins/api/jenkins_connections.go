package api

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/utils"
)

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jenkinsSource, err := GetSourceFromEnv()
	if err != nil {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.ReadError}}, nil
	}
	if jenkinsSource.Username == "" || jenkinsSource.Password == "" || jenkinsSource.Endpoint == "" {
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.UnsetConnectionError}}, nil
	}
	encodedToken := utils.GetEncodedToken(jenkinsSource.Username, jenkinsSource.Password)
	apiClient := &core.ApiClient{}
	apiClient.Setup(
		jenkinsSource.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", encodedToken),
		},
		10*time.Second,
		3,
	)
	res, err := apiClient.Get("", nil, nil)
	if err != nil || res.StatusCode != 200 {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.InvalidConnectionError}}, nil
	}
	return &core.ApiResourceOutput{Body: core.TestResult{Success: true, Message: ""}}, nil
}
