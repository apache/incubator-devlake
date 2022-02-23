package api

import (
	"net/http"
	"net/url"
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks"
	"github.com/merico-dev/lake/utils"
	"github.com/mitchellh/mapstructure"
)

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint"`
	Auth     string `json:"auth"`
	Proxy    string `json:"proxy"`
}

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	ValidationResult := core.ValidateParams(input, []string{"endpoint", "auth"})
	if !ValidationResult.Success {
		return &core.ApiResourceOutput{Body: ValidationResult}, nil
	}
	var params TestConnectionRequest
	err := mapstructure.Decode(input.Body, &params)
	if err != nil {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.InvalidParams}}, nil
	}

	parsedUrl, err := url.Parse(params.Endpoint)
	if err != nil {
		// parsed error
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.InvalidEndpointError}}, nil
	}
	if parsedUrl.Scheme == "" {
		// no schema
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.SchemaIsRequired}}, nil
	}
	err = utils.CheckDNS(parsedUrl.Hostname())
	if err != nil {
		// ip not found
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.DNSResolveFailedError}}, nil
	}
	port, err := utils.ResolvePort(parsedUrl.Port(), parsedUrl.Scheme)
	if err != nil {
		// resolve port failed
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.InvalidSchema}}, nil
	}
	err = utils.CheckNetwork(parsedUrl.Hostname(), port, time.Duration(2)*time.Second)
	if err != nil {
		// connect failed
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: core.NetworkConnectError}}, nil
	}

	jiraApiClient := tasks.NewJiraApiClient(params.Endpoint, params.Auth, params.Proxy, nil)
	jiraApiClient.SetTimeout(2 * time.Second)

	serverInfo, statusCode, err := jiraApiClient.GetJiraServerInfo()
	if statusCode == http.StatusNotFound {
		// failed to get jira version
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: InvaildJiraApi}}, nil
	}
	if statusCode == http.StatusUnauthorized {
		// NOTICE: jira api will check your token if you provided it, even the api can be accessed anonymously
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: InvalidAuthInfo}}, nil
	}
	if err != nil {
		logger.Error("Error: ", err)
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: err.Error()}}, nil
	}

	if serverInfo.DeploymentType != models.DeploymentCloud {
		// unsupported jira version
		// FIXME: remove it when jira server is supported
		return &core.ApiResourceOutput{Body: core.TestResult{Success: false, Message: InvalidJiraVersion}}, nil
	}

	return &core.ApiResourceOutput{Body: core.TestResult{Success: true, Message: ""}}, nil
}

const InvaildJiraApi = "Failed to request jira version api"
const InvalidJiraVersion = "Unsupported jira server, only support jira cloud version now"
const InvalidAuthInfo = "Authentication failed, please check your Basic Auth Token"
