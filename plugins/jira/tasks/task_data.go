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

	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type JiraOptions struct {
	ConnectionId    uint64   `json:"connectionId"`
	BoardId         uint64   `json:"boardId"`
	Tasks           []string `json:"tasks,omitempty"`
	Since           string
	IssueExtraction struct {
		RequirementTypeMapping []string
		BugTypeMapping         []string
		IncidentTypeMapping    []string
		//EpicKeyField           string
		StoryPointField string
	}
}

type JiraTaskData struct {
	Options        *JiraOptions
	ApiClient      *helper.ApiAsyncClient
	Connection     *models.JiraConnection
	Since          *time.Time
	JiraServerInfo models.JiraServerInfo
}
