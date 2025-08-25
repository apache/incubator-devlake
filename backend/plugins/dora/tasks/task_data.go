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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type DoraApiParams struct {
	ProjectName string
}

type DoraOptions struct {
	// options for dora plugin, required when activate dora
	ProjectName string `json:"projectName"`
	RepoUrl     string `json:"repoUrl"` // optional, for locating repo
	// Time after which the data represents the real history, like github cloud is 2015-01-01, or the project start date
	TimeAfter *time.Time `json:"timeAfter"`

	// Specify the issue statuses that indicate 'In Progress'. Can be multiple comma-separated values.
	// e.g., "In Progress,In Development,In Review"
	InProgressStatus string `json:"inProgressStatus"`

	// Specify the issue statuses that indicate 'Done'. Can be multiple comma-separated values.
	// e.g., "Done,Closed,Resolved"
	DoneStatus string `json:"doneStatus"`

	// --- Keep original fields needed by other DORA tasks ---
	ScopeConfigId   uint64 `json:"scopeConfigId"`
	ScopeConfigName string `json:"scopeConfigName"`
}

func (o *DoraOptions) GetInProgressStatuses() []string {
	if o.InProgressStatus == "" {
		return []string{}
	}
	return strings.Split(o.InProgressStatus, ",")
}

func (o *DoraOptions) GetDoneStatuses() []string {
	if o.DoneStatus == "" {
		return []string{}
	}
	return strings.Split(o.DoneStatus, ",")
}

type DoraTaskData struct {
	Options   *DoraOptions
	TimeAfter *time.Time
	RepoUrl   string `json:"repoUrl"` // optional, for locating repo

	// --- Keep fields needed by other DORA tasks ---
	ScopeId                         string `json:"scopeId"`
	ScopeName                       string `json:"scopeName"`
	DisableIssueToIncidentGenerator bool   `json:"disableIssueToIncidentGenerator"`

	RemoteArgs interface{} // Temporary workaround to let it compile
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*DoraOptions, errors.Error) {
	var op DoraOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	// find the scope config if owner/repo is not specified
	if op.ProjectName == "" {
		return nil, errors.BadInput.New("projectName is required for dora plugin")
	}
	return &op, nil
}
