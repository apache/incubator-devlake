package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/lake/plugins/core"
)

type ApiUserPublicEmailResponse []PublicEmail

// Using Public Email because it requires authentication, and it is public information anyway.
// We're not using email information for anything here.
type PublicEmail struct {
	Email      string
	Primary    bool
	Verified   bool
	Visibility string
}

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required,url"`
	Auth     string `json:"auth" validate:"required"`
	Proxy    string `json:"proxy"`
}

var vld = validator.New()

/*
POST /plugins/github/test
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
	tokens := strings.Split(params.Auth, ",")

	// verify multiple token in parallel
	// PLEASE NOTE: This works because GitHub API Client rotates tokens on each request
	results := make(chan error)
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		i := i
		go func() {
			apiClient, err := helper.NewApiClient(
				params.Endpoint,
				map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", token),
				},
				3*time.Second,
				params.Proxy,
				nil,
			)
			if err != nil {
				results <- fmt.Errorf("verify token failed for #%v %s %w", i, token, err)
				return
			}
			res, err := apiClient.Get("user/public_emails", nil, nil)
			if err != nil {
				results <- fmt.Errorf("verify token failed for #%v %s %w", i, token, err)
				return
			}
			githubApiResponse := &ApiUserPublicEmailResponse{}
			err = helper.UnmarshalResponse(res, githubApiResponse)
			if err != nil {
				results <- fmt.Errorf("verify token failed for #%v %s %w", i, token, err)
			} else {
				results <- nil
			}
		}()
	}

	// collect verification results
	msgs := make([]string, 0)
	i := 0
	for err := range results {
		if err != nil {
			msgs = append(msgs, err.Error())
		}
		i++
		if i == len(tokens) {
			close(results)
		}
	}

	// output
	return &core.ApiResourceOutput{
		Body: core.TestResult{
			Success: len(msgs) == 0,
			Message: strings.Join(msgs, "\n"),
		},
	}, nil
}
