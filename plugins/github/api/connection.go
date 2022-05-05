package api

import (
	"fmt"
	"github.com/merico-dev/lake/config"
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
	Endpoint string `mapstructure:"GITHUB_ENDPOINT" validate:"required"`
	Auth     string `mapstructure:"GITHUB_AUTH" validate:"required"`
	Proxy    string `mapstructure:"GITHUB_PROXY"`

	PrType               string `mapstructure:"GITHUB_PR_TYPE"`
	PrComponent          string `mapstructure:"GITHUB_PR_COMPONENT"`
	IssueSeverity        string `mapstructure:"GITHUB_ISSUE_SEVERITY"`
	IssuePriority        string `mapstructure:"GITHUB_ISSUE_PRIORITY"`
	IssueComponent       string `mapstructure:"GITHUB_ISSUE_COMPONENT"`
	IssueTypeBug         string `mapstructure:"GITHUB_ISSUE_TYPE_BUG"`
	IssueTypeIncident    string `mapstructure:"GITHUB_ISSUE_TYPE_INCIDENT"`
	IssueTypeRequirement string `mapstructure:"GITHUB_ISSUE_TYPE_REQUIREMENT"`
}

// This object conforms to what the frontend currently expects.
type GithubResponse struct {
	Name string
	ID   int

	GithubConnection
}

/*
PUT /plugins/github/connections/:connectionId
*/
func PutConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	githubConnection := GithubConnection{}
	err := mapstructure.Decode(input.Body, &githubConnection)
	if err != nil {
		return nil, err
	}

	err = config.SetStruct(githubConnection, "mapstructure")
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: "Success"}, nil
}

/*
GET /plugins/github/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 Connection (FROM ENV) until multi-connection is developed.
	githubConnections, err := GetConnectionFromEnv()
	response := []GithubResponse{*githubConnections}
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: response}, nil
}

/*
GET /plugins/github/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-connection is developed.)
	githubConnections, err := GetConnectionFromEnv()
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: githubConnections}, nil
}

func GetConnectionFromEnv() (*GithubResponse, error) {
	v := config.GetConfig()
	var githubConnection GithubConnection
	err := v.Unmarshal(&githubConnection)
	if err != nil {
		return nil, err
	}

	return &GithubResponse{
		Name:             "Github",
		ID:               1,
		GithubConnection: githubConnection,
	}, nil
}
