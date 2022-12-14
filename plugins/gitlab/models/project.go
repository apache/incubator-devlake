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
	"time"

	"github.com/apache/incubator-devlake/models/common"
)

type GitlabProject struct {
	ConnectionId            uint64 `json:"connectionId" mapstructure:"connectionId" gorm:"primaryKey"`
	TransformationRuleId    uint64 `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId"`
	GitlabId                int    `json:"gitlabId" mapstructure:"gitlabId" gorm:"primaryKey"`
	Name                    string `json:"name" mapstructure:"name" gorm:"type:varchar(255)"`
	Description             string `json:"description" mapstructure:"description"`
	DefaultBranch           string `json:"defaultBranch" mapstructure:"defaultBranch" gorm:"type:varchar(255)"`
	PathWithNamespace       string `json:"pathWithNamespace" mapstructure:"pathWithNamespace" gorm:"type:varchar(255)"`
	WebUrl                  string `json:"webUrl" mapstructure:"webUrl" gorm:"type:varchar(255)"`
	CreatorId               int    `json:"creatorId" mapstructure:"creatorId"`
	Visibility              string `json:"visibility" mapstructure:"visibility" gorm:"type:varchar(255)"`
	OpenIssuesCount         int    `json:"openIssuesCount" mapstructure:"openIssuesCount"`
	StarCount               int    `json:"starCount" mapstructure:"StarCount"`
	ForkedFromProjectId     int    `json:"forkedFromProjectId" mapstructure:"forkedFromProjectId"`
	ForkedFromProjectWebUrl string `json:"forkedFromProjectWebUrl" mapstructure:"forkedFromProjectWebUrl" gorm:"type:varchar(255)"`
	HttpUrlToRepo           string `json:"http_url_to_repo" gorm:"varchar(255)"`

	CreatedDate      time.Time  `json:"createdDate" mapstructure:"-"`
	UpdatedDate      *time.Time `json:"updatedDate" mapstructure:"-"`
	common.NoPKModel `json:"-" mapstructure:"-"`
}

func (GitlabProject) TableName() string {
	return "_tool_gitlab_projects"
}
