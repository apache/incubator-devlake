/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package models

type GithubConnection struct {
	Endpoint string `mapstructure:"endpoint" validate:"required" env:"GITHUB_ENDPOINT" json:"endpoint"`
	Auth     string `mapstructure:"auth" validate:"required" env:"GITHUB_AUTH" json:"auth"`
	Proxy    string `mapstructure:"proxy" env:"GITHUB_PROXY" json:"proxy"`

	Config `mapstructure:",squash"`
}

type Config struct {
	PrType               string `mapstructure:"prType" env:"GITHUB_PR_TYPE" json:"prType"`
	PrComponent          string `mapstructure:"prComponent" env:"GITHUB_PR_COMPONENT" json:"prComponent"`
	IssueSeverity        string `mapstructure:"issueSeverity" env:"GITHUB_ISSUE_SEVERITY" json:"issueSeverity"`
	IssuePriority        string `mapstructure:"issuePriority" env:"GITHUB_ISSUE_PRIORITY" json:"issuePriority"`
	IssueComponent       string `mapstructure:"issueComponent" env:"GITHUB_ISSUE_COMPONENT" json:"issueComponent"`
	IssueTypeBug         string `mapstructure:"issueTypeBug" env:"GITHUB_ISSUE_TYPE_BUG" json:"issueTypeBug"`
	IssueTypeIncident    string `mapstructure:"issueTypeIncident" env:"GITHUB_ISSUE_TYPE_INCIDENT" json:"issueTypeIncident"`
	IssueTypeRequirement string `mapstructure:"issueTypeRequirement" env:"GITHUB_ISSUE_TYPE_REQUIREMENT" json:"issueTypeRequirement"`
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
