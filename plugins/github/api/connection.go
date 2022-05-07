package api

import (
	"fmt"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/github/models"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/lake/plugins/core"
)

type ApiUserPublicEmailResponse []models.PublicEmail

var vld = validator.New()

/*
POST /plugins/github/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// process input
	var params models.TestConnectionRequest
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

	if len(msgs) > 0 {
		return nil, fmt.Errorf(strings.Join(msgs, "\n"))
	}

	// output
	return nil, nil
}

/*
PATCH /plugins/github/connections/:connectionId
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	v := config.GetConfig()
	connection := &models.GithubConnection{}
	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	// update from request and save to .env
	err = helper.DecodeStruct(v, connection, input.Body, "env")
	if err != nil {
		return nil, err
	}
	err = v.WriteConfig()
	if err != nil {
		return nil, err
	}
	response := models.GithubResponse{
		GithubConnection: *connection,
		Name:             "Github",
		ID:               1,
	}
	return &core.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

/*
GET /plugins/github/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-connection is developed.
	v := config.GetConfig()
	connection := &models.GithubConnection{}

	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	response := models.GithubResponse{
		GithubConnection: *connection,
		Name:             "Github",
		ID:               1,
	}

	return &core.ApiResourceOutput{Body: []models.GithubResponse{response}}, nil
}

/*
GET /plugins/github/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-connection is developed.)
	v := config.GetConfig()
	connection := &models.GithubConnection{}
	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	response := &models.GithubResponse{
		GithubConnection: *connection,
		Name:             "Github",
		ID:               1,
	}
	return &core.ApiResourceOutput{Body: response}, nil
}
