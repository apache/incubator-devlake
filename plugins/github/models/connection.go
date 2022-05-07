package models

type GithubConnection struct {
	Endpoint string `mapstructure:"endpoint" validate:"required" env:"GITHUB_ENDPOINT" json:"endpoint"`
	Auth     string `mapstructure:"auth" validate:"required" env:"GITHUB_AUTH" json:"auth"`
	Proxy    string `mapstructure:"proxy" env:"GITHUB_PROXY" json:"proxy"`

	Config
}

type Config struct {
	PrType               string `mapstructure:"prType,squash" env:"GITHUB_PR_TYPE" json:"prType"`
	PrComponent          string `mapstructure:"prComponent,squash" env:"GITHUB_PR_COMPONENT" json:"prComponent"`
	IssueSeverity        string `mapstructure:"issueSeverity,squash" env:"GITHUB_ISSUE_SEVERITY" json:"issueSeverity"`
	IssuePriority        string `mapstructure:"issuePriority,squash" env:"GITHUB_ISSUE_PRIORITY" json:"issuePriority"`
	IssueComponent       string `mapstructure:"issueComponent,squash" env:"GITHUB_ISSUE_COMPONENT" json:"issueComponent"`
	IssueTypeBug         string `mapstructure:"issueTypeBug,squash" env:"GITHUB_ISSUE_TYPE_BUG" json:"issueTypeBug"`
	IssueTypeIncident    string `mapstructure:"typeIncident,squash" env:"GITHUB_ISSUE_TYPE_INCIDENT" json:"typeIncident"`
	IssueTypeRequirement string `mapstructure:"issueTypeRequirement,squash" env:"GITHUB_ISSUE_TYPE_REQUIREMENT" json:"issueTypeRequirement"`
}

// This object conforms to what the frontend currently expects.
type GithubResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`

	GithubConnection
}

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
