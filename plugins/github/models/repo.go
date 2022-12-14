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
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

type GithubRepo struct {
	ConnectionId         uint64     `json:"connectionId" gorm:"primaryKey" mapstructure:"connectionId,omitempty"`
	GithubId             int        `json:"githubId" gorm:"primaryKey" mapstructure:"githubId"`
	Name                 string     `json:"name" gorm:"type:varchar(255)" mapstructure:"name,omitempty"`
	HTMLUrl              string     `json:"HTMLUrl" gorm:"type:varchar(255)" mapstructure:"HTMLUrl,omitempty"`
	Description          string     `json:"description" mapstructure:"description,omitempty"`
	TransformationRuleId uint64     `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId,omitempty"`
	OwnerId              int        `json:"ownerId" mapstructure:"ownerId,omitempty"`
	Language             string     `json:"language" gorm:"type:varchar(255)" mapstructure:"language,omitempty"`
	ParentGithubId       int        `json:"parentId" mapstructure:"parentGithubId,omitempty"`
	ParentHTMLUrl        string     `json:"parentHtmlUrl" mapstructure:"parentHtmlUrl,omitempty"`
	CloneUrl             string     `json:"cloneUrl" gorm:"type:varchar(255)" mapstructure:"cloneUrl,omitempty"`
	CreatedDate          *time.Time `json:"createdDate" mapstructure:"-"`
	UpdatedDate          *time.Time `json:"updatedDate" mapstructure:"-"`
	common.NoPKModel     `json:"-" mapstructure:"-"`
}

func (GithubRepo) TableName() string {
	return "_tool_github_repos"
}
