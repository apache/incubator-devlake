package models

type DeploymentType string
type Locale string

const DeploymentCloud DeploymentType = "Cloud"
const DeploymentServer DeploymentType = "Server"
const LocaleEnUS Locale = "en_US"

type JiraServerInfo struct {
	BaseURL        string         `json:"baseUrl"`
	BuildDate      string         `json:"buildDate"`
	BuildNumber    int            `json:"buildNumber"`
	DeploymentType DeploymentType `json:"deploymentType"`
	ScmInfo        string         `json:"ScmInfo"`
	ServerTime     string         `json:"serverTime"`
	ServerTitle    string         `json:"serverTitle"`
	Version        string         `json:"version"`
	VersionNumbers []int          `json:"versionNumbers"`
}

type ApiMyselfResponse struct {
	AccountId   string
	DisplayName string
}

func (JiraServerInfo) TableName() string{
	return "_tool_jira_server_infos"
}

func (ApiMyselfResponse) TableName() string{
	return "_tool_api_myself_responses"
}

