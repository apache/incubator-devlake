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

import (
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type GiteeAccessToken helper.AccessToken

// SetupAuthentication sets up the HTTP Request Authentication
func (gat GiteeAccessToken) SetupAuthentication(req *http.Request) errors.Error {
	query := req.URL.Query()
	query.Set("access_token", gat.Token)
	req.URL.RawQuery = query.Encode()
	return nil
}

// GiteeConn holds the essential information to connect to the Gitee API
type GiteeConn struct {
	helper.RestConnection `mapstructure:",squash"`
	GiteeAccessToken      `mapstructure:",squash"`
}

// GiteeConnection holds GiteeConn plus ID/Name for database storage
type GiteeConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	GiteeConn             `mapstructure:",squash"`
}

type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

type TransformationRules struct {
	PrType               string `mapstructure:"prType" env:"GITEE_PR_TYPE" json:"prType"`
	PrComponent          string `mapstructure:"prComponent" env:"GITEE_PR_COMPONENT" json:"prComponent"`
	PrBodyClosePattern   string `mapstructure:"prBodyClosePattern" json:"prBodyClosePattern"`
	IssueSeverity        string `mapstructure:"issueSeverity" env:"GITEE_ISSUE_SEVERITY" json:"issueSeverity"`
	IssuePriority        string `mapstructure:"issuePriority" env:"GITEE_ISSUE_PRIORITY" json:"issuePriority"`
	IssueComponent       string `mapstructure:"issueComponent" env:"GITEE_ISSUE_COMPONENT" json:"issueComponent"`
	IssueTypeBug         string `mapstructure:"issueTypeBug" env:"GITEE_ISSUE_TYPE_BUG" json:"issueTypeBug"`
	IssueTypeIncident    string `mapstructure:"issueTypeIncident" env:"GITEE_ISSUE_TYPE_INCIDENT" json:"issueTypeIncident"`
	IssueTypeRequirement string `mapstructure:"issueTypeRequirement" env:"GITEE_ISSUE_TYPE_REQUIREMENT" json:"issueTypeRequirement"`
	DeploymentPattern    string `mapstructure:"deploymentPattern" json:"deploymentPattern"`
}

func (GiteeConnection) TableName() string {
	return "_tool_gitee_connections"
}
