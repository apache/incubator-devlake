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

