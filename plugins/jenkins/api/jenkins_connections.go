package api

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/utils"
)

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	ValidationResult := core.ValidateParams(input, []string{"username", "password", "endpoint"})
	if !ValidationResult.Success {
		return &core.ApiResourceOutput{Body: ValidationResult}, nil
	}

	encodedToken := utils.GetEncodedToken(input.Query.Get("username"), input.Query.Get("password"))
	apiClient := &core.ApiClient{}
	apiClient.Setup(
		input.Query.Get("endpoint"),
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
