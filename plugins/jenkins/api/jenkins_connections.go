package api

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/utils"
)

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jenkinsSource, err := GetSourceFromEnv()
	if err != nil {
		return nil, err
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
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Your connection configuration is invalid.")
	}
	return &core.ApiResourceOutput{Body: true}, nil
}
