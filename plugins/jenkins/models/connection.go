package models

// This object conforms to what the frontend currently sends.
type JenkinsConnection struct {
	Endpoint string `mapstructure:"endpoint" validate:"required" env:"JENKINS_ENDPOINT" json:"endpoint"`
	Username string `mapstructure:"username" validate:"required" env:"JENKINS_USERNAME" json:"username"`
	Password string `mapstructure:"password" validate:"required" env:"JENKINS_PASSWORD" json:"password"`
	Proxy    string `mapstructure:"proxy" env:"JENKINS_PROXY" json:"proxy"`
}

type JenkinsResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	JenkinsConnection
}

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Proxy    string `json:"proxy"`
}
