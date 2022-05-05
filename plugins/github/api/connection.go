package api

import (
	"fmt"
	"github.com/merico-dev/lake/config"
	"net/http"
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

	if len(msgs) > 0 {
		return nil, fmt.Errorf(strings.Join(msgs, "\n"))
	}

	// output
	return nil, nil
}

// This object conforms to what the frontend currently sends.
type GithubConnection struct {
	Endpoint string `mapstructure:"endpoint" validate:"required" env:"GITHUB_ENDPOINT" json:"endpoint"`
	Auth     string `mapstructure:"auth" validate:"required" env:"GITHUB_AUTH" json:"auth"`
	Proxy    string `mapstructure:"proxy" env:"GITHUB_PROXY" json:"proxy"`

	PrType               string `mapstructure:"prType" env:"GITHUB_PR_TYPE" json:"prType"`
	PrComponent          string `mapstructure:"prComponent" env:"GITHUB_PR_COMPONENT" json:"prComponent"`
	IssueSeverity        string `mapstructure:"issueSeverity" env:"GITHUB_ISSUE_SEVERITY" json:"issueSeverity"`
	IssuePriority        string `mapstructure:"issuePriority" env:"GITHUB_ISSUE_PRIORITY" json:"issuePriority"`
	IssueComponent       string `mapstructure:"issueComponent" env:"GITHUB_ISSUE_COMPONENT" json:"issueComponent"`
	IssueTypeBug         string `mapstructure:"issueTypeBug" env:"GITHUB_ISSUE_TYPE_BUG" json:"issueTypeBug"`
	IssueTypeIncident    string `mapstructure:"typeIncident" env:"GITHUB_ISSUE_TYPE_INCIDENT" json:"typeIncident"`
	IssueTypeRequirement string `mapstructure:"issueTypeRequirement" env:"GITHUB_ISSUE_TYPE_REQUIREMENT" json:"issueTypeRequirement"`
}

// This object conforms to what the frontend currently expects.
type GithubResponse struct {
	Name string
	ID   int

	GithubConnection
}

/*
PATCH /plugins/github/connections/:connectionId
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &GithubConnection{}, "env")
	if err != nil {
		return nil, err
	}
	// update from request and save to .env
	err = helper.SaveToConifgWithMap(v, connection, input.Body, "env")
	if err != nil {
		return nil, err
	}

	response := GithubResponse{
		GithubConnection: *connection.(*GithubConnection),
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
	connection, err := helper.LoadFromConfig(v, &GithubConnection{}, "env")
	if err != nil {
		return nil, err
	}
	response := GithubResponse{
		GithubConnection: *connection.(*GithubConnection),
		Name:             "Github",
		ID:               1,
	}

	return &core.ApiResourceOutput{Body: []GithubResponse{response}}, nil
}

/*
GET /plugins/github/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-connection is developed.)
	v := config.GetConfig()
	connection, err := helper.LoadFromConfig(v, &GithubConnection{}, "env")
	if err != nil {
		return nil, err
	}
	response := &GithubResponse{
		GithubConnection: *connection.(*GithubConnection),
		Name:             "Github",
		ID:               1,
	}
	return &core.ApiResourceOutput{Body: response}, nil
}
