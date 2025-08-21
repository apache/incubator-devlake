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
	"github.com/apache/incubator-devlake/core/plugin"

	"github.com/apache/incubator-devlake/core/models/common"
)

var _ plugin.ToolLayerScope = (*AzuredevopsProject)(nil)

type AzuredevopsProject struct {
	common.Scope `mapstructure:",squash" gorm:"embedded"`

	AzuredevopsId  string             `json:"id"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	Url            string             `json:"url"`
	State          string             `json:"state"`
	Revision       int                `json:"revision"`
	Visibility     string             `json:"visibility"`
	LastUpdateTime common.Iso8601Time `json:"lastUpdateTime"`
}

func (a AzuredevopsProject) TableName() string {
	return "_tool_azuredevops_go_projects"
}

func (a AzuredevopsProject) ScopeId() string {
	return a.AzuredevopsId
}

func (a AzuredevopsProject) ScopeName() string {
	return a.Name
}

func (a AzuredevopsProject) ScopeFullName() string {
	return a.Name
}

func (a AzuredevopsProject) ScopeParams() interface{} {
	//TODO implement me
	panic("implement me")
}
