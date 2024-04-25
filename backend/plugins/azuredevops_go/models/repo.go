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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.ToolLayerScope = (*AzuredevopsRepo)(nil)

type AzuredevopsRepo struct {
	common.Scope  `mapstructure:",squash"`
	AzureDevOpsPK `mapstructure:",squash"`

	Id         string `json:"id" validate:"required" mapstructure:"id" gorm:"primaryKey"`
	Type       string `json:"type" validate:"required" mapstructure:"type"`
	Name       string `json:"name" mapstructure:"name,omitempty"`
	Url        string `json:"url" mapstructure:"url,omitempty"`
	RemoteUrl  string `json:"remoteUrl"`
	ExternalId string
	IsFork     bool
	IsPrivate  bool
}

func (repo AzuredevopsRepo) ScopeId() string {
	return repo.Id
}

func (repo AzuredevopsRepo) ScopeName() string {
	return repo.Name
}

func (repo AzuredevopsRepo) ScopeFullName() string {
	return repo.Name
}

func (repo AzuredevopsRepo) ScopeParams() interface{} {
	return &AzuredevopsApiParams{
		ConnectionId: repo.ConnectionId,
		Name:         repo.Name,
	}
}

func (AzuredevopsRepo) TableName() string {
	return "_tool_azuredevops_go_repos"
}

type AzuredevopsApiParams struct {
	ConnectionId uint64
	Name         string
}
