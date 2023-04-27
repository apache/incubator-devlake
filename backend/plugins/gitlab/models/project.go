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
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"

	"github.com/apache/incubator-devlake/core/models/common"
)

var _ plugin.ToolLayerScope = (*GitlabProject)(nil)

type GitlabProject struct {
	ConnectionId            uint64 `json:"connectionId" mapstructure:"connectionId" validate:"required" gorm:"primaryKey"`
	TransformationRuleId    uint64 `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId"`
	GitlabId                int    `json:"gitlabId" mapstructure:"gitlabId" validate:"required" gorm:"primaryKey"`
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
	HttpUrlToRepo           string `json:"httpUrlToRepo" gorm:"type:varchar(255)"`

	CreatedDate      *time.Time `json:"createdDate" mapstructure:"-"`
	UpdatedDate      *time.Time `json:"updatedDate" mapstructure:"-"`
	common.NoPKModel `json:"-" mapstructure:"-"`
}

func (GitlabProject) TableName() string {
	return "_tool_gitlab_projects"
}

func (p GitlabProject) ScopeId() string {
	return strconv.Itoa(p.GitlabId)
}

func (p GitlabProject) ScopeName() string {
	return p.Name
}

// Convert the API response to our DB model instance
func (gitlabApiProject GitlabApiProject) ConvertApiScope() plugin.ToolLayerScope {
	p := &GitlabProject{}
	p.GitlabId = gitlabApiProject.GitlabId
	p.Name = gitlabApiProject.Name
	p.Description = gitlabApiProject.Description
	p.DefaultBranch = gitlabApiProject.DefaultBranch
	p.CreatorId = gitlabApiProject.CreatorId
	p.PathWithNamespace = gitlabApiProject.PathWithNamespace
	p.WebUrl = gitlabApiProject.WebUrl
	p.HttpUrlToRepo = gitlabApiProject.HttpUrlToRepo
	p.Visibility = gitlabApiProject.Visibility
	p.OpenIssuesCount = gitlabApiProject.OpenIssuesCount
	p.StarCount = gitlabApiProject.StarCount
	p.CreatedDate = gitlabApiProject.CreatedAt.ToNullableTime()
	p.UpdatedDate = helper.Iso8601TimeToTime(gitlabApiProject.LastActivityAt)
	if gitlabApiProject.ForkedFromProject != nil {
		p.ForkedFromProjectId = gitlabApiProject.ForkedFromProject.GitlabId
		p.ForkedFromProjectWebUrl = gitlabApiProject.ForkedFromProject.WebUrl
	}
	// this might happen when GitlabConnection.SearchScopes
	if len(p.Name) > len(p.PathWithNamespace) {
		p.Name, p.PathWithNamespace = p.PathWithNamespace, p.Name
	}
	return p
}

type GitlabApiProject struct {
	GitlabId          int    `json:"id"`
	Name              string `josn:"name"`
	Description       string `json:"description"`
	DefaultBranch     string `json:"default_branch"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebUrl            string `json:"web_url"`
	CreatorId         int
	Visibility        string              `json:"visibility"`
	OpenIssuesCount   int                 `json:"open_issues_count"`
	StarCount         int                 `json:"star_count"`
	ForkedFromProject *GitlabApiProject   `json:"forked_from_project"`
	CreatedAt         helper.Iso8601Time  `json:"created_at"`
	LastActivityAt    *helper.Iso8601Time `json:"last_activity_at"`
	HttpUrlToRepo     string              `json:"http_url_to_repo"`
}

type GroupResponse struct {
	Id          int    `json:"id" group:"id"`
	WebUrl      string `json:"web_url"`
	Name        string `json:"name" group:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	FullName    string `json:"full_name"`
	FullPath    string `json:"full_path"`
	ParentId    *int   `json:"parent_id"`
}

func (p GroupResponse) GroupId() string {
	return "group:" + strconv.Itoa(p.Id)
}

func (p GroupResponse) GroupName() string {
	return p.Name
}
