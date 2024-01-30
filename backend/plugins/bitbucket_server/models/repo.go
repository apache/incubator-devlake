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
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.ToolLayerScope = (*BitbucketServerRepo)(nil)
var _ plugin.ApiGroup = (*ProjectItem)(nil)
var _ plugin.ApiScope = (*BitbucketServerApiRepo)(nil)

type BitbucketServerRepo struct {
	common.Scope `mapstructure:",squash"`
	BitbucketId  string     `json:"bitbucketId" gorm:"primaryKey;type:varchar(255)" validate:"required" mapstructure:"bitbucketId"`
	Name         string     `json:"name" gorm:"type:varchar(255)" mapstructure:"name,omitempty"`
	HTMLUrl      string     `json:"HTMLUrl" gorm:"type:varchar(255)" mapstructure:"HTMLUrl,omitempty"`
	Description  string     `json:"description" mapstructure:"description,omitempty"`
	CloneUrl     string     `json:"cloneUrl" gorm:"type:varchar(255)" mapstructure:"cloneUrl,omitempty"`
	CreatedDate  *time.Time `json:"createdDate" mapstructure:"-"`
	UpdatedDate  *time.Time `json:"updatedDate" mapstructure:"-"`
}

func (BitbucketServerRepo) TableName() string {
	return "_tool_bitbucket_server_repos"
}

func (p BitbucketServerRepo) ScopeId() string {
	return p.BitbucketId
}

func (p BitbucketServerRepo) ScopeName() string {
	return p.Name
}

func (p BitbucketServerRepo) ScopeFullName() string {
	return p.BitbucketId
}

func (p BitbucketServerRepo) ScopeParams() interface{} {
	return &BitbucketServerApiParams{
		ConnectionId: p.ConnectionId,
		FullName:     p.BitbucketId,
	}
}

type BitbucketServerApiRepo struct {
	Id            int32  `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	HierarchyId   string `json:"hierarchyId"`
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Description   string `json:"description"`
	Public        bool   `json:"public"`
	Archived      bool   `json:"archived"`
	// CreatedAt     *time.Time  `json:"created_on"`
	// UpdatedAt     *time.Time  `json:"updated_on"`
	Project ProjectItem `json:"project"`
	Links   struct {
		Clone []struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"clone"`
		Self []struct {
			Href string `json:"href"`
		}
	} `json:"links"`
}

func (b BitbucketServerApiRepo) ConvertApiScope() plugin.ToolLayerScope {
	scope := &BitbucketServerRepo{}
	scope.BitbucketId = fmt.Sprintf("%s/repos/%s", b.Project.Key, b.Slug)
	scope.Description = b.Description
	scope.Name = b.Name

	if len(b.Links.Self) > 0 {
		scope.HTMLUrl = b.Links.Self[0].Href
	}

	scope.CloneUrl = ""
	for _, u := range b.Links.Clone {
		if u.Name == "http" {
			scope.CloneUrl = u.Href
		}
	}
	return scope
}

type PaginationResponse[T any] struct {
	Start         int  `json:"start"`
	Limit         int  `json:"limit"`
	Size          int  `json:"size"`
	IsLastPage    bool `json:"isLastPage"`
	NextPageStart *int `json:"nextPageStart"`
	Values        []T  `json:"values"`
}

type ProjectsResponse PaginationResponse[ProjectItem]
type ReposResponse PaginationResponse[BitbucketServerApiRepo]

type ProjectItem struct {
	Key    string `json:"key" group:"id"`
	Name   string `json:"name" group:"name"`
	Public bool   `json:"public"`
	Type   string `json:"type"`
	Id     int32  `json:"id"`
}

func (p ProjectItem) GroupId() string {
	return p.Key
}

func (p ProjectItem) GroupName() string {
	return p.Name
}

type BitbucketServerApiParams struct {
	ConnectionId uint64
	FullName     string
}
