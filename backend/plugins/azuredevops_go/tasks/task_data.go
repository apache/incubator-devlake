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

package tasks

import (
	"time"

	"github.com/apache/incubator-devlake/core/errors"

	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
)

type AzuredevopsOptions struct {
	ConnectionId   uint64 `json:"connectionId" mapstructure:"connectionId,omitempty"`
	ProjectId      string `json:"projectId" mapstructure:"projectId,omitempty"`
	OrganizationId string `json:"organizationId" mapstructure:"organizationId,omitempty"`
	RepositoryId   string `json:"repositoryId"  mapstructure:"repositoryId,omitempty"`
	RepositoryType string `json:"repositoryType"  mapstructure:"repositoryType,omitempty"`
	ExternalId     string `json:"externalId"  mapstructure:"externalId,omitempty"`

	ScopeConfigId uint64                         `json:"scopeConfigId" mapstructure:"scopeConfigId,omitempty"`
	TimeAfter     string                         `json:"timeAfter" mapstructure:"timeAfter,omitempty"`
	ScopeConfig   *models.AzuredevopsScopeConfig `mapstructure:"scopeConfig,omitempty" json:"scopeConfig"`
}

type AzuredevopsTaskData struct {
	Options       *AzuredevopsOptions
	ApiClient     *helper.ApiAsyncClient
	TimeAfter     *time.Time
	RegexEnricher *helper.RegexEnricher
}

func DecodeTaskOptions(options map[string]interface{}) (*AzuredevopsOptions, errors.Error) {
	var op AzuredevopsOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

type AzuredevopsParams struct {
	OrganizationId string
	RepositoryId   string
	ProjectId      string
}

func (p *AzuredevopsOptions) GetParams() any {
	return AzuredevopsParams{
		OrganizationId: p.OrganizationId,
		ProjectId:      p.ProjectId,
		RepositoryId:   p.RepositoryId,
	}
}
