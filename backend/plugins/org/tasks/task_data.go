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

import "github.com/apache/incubator-devlake/core/plugin"

type Options struct {
	ConnectionId    uint64           `json:"connectionId"`
	ProjectMappings []ProjectMapping `json:"projectMappings"`
}

// ProjectMapping represents the relations between project and scopes
type ProjectMapping struct {
	ProjectName string  `json:"projectName"`
	Scopes      []Scope `json:"scopes"`
}

// Scope represents a scope by specifies the table and id
type Scope struct {
	Table string `json:"table"`
	RowID string `json:"rowId"`
}

// NewProjectMapping is the construct function of ProjectMapping
func NewProjectMapping(projectName string, pluginScopes []plugin.Scope) ProjectMapping {
	var scopes []Scope
	for _, ps := range pluginScopes {
		scopes = append(scopes, Scope{
			Table: ps.TableName(),
			RowID: ps.ScopeId(),
		})
	}
	return ProjectMapping{
		ProjectName: projectName,
		Scopes:      scopes,
	}
}

type TaskData struct {
	Options *Options
}
type Params struct {
	ConnectionId uint64
}
