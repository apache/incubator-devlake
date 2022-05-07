package models

// This object conforms to what the frontend currently sends.
type GitlabConnection struct {
	Endpoint string `mapstructure:"endpoint" validate:"required" env:"GITLAB_ENDPOINT" json:"endpoint"`
	Auth     string `mapstructure:"auth" validate:"required" env:"GITLAB_AUTH"  json:"auth"`
	Proxy    string `mapstructure:"proxy" env:"GITLAB_PROXY" json:"proxy"`
}

// This object conforms to what the frontend currently expects.
type GitlabResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	GitlabConnection
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required"`
	Auth     string `json:"auth" validate:"required"`
	Proxy    string `json:"proxy"`
}
