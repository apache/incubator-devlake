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

package apiv2models

import (
	"github.com/apache/incubator-devlake/plugins/helper"
	"time"

	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type Worklog struct {
	Self             string             `json:"self"`
	Author           *User              `json:"author"`
	UpdateAuthor     *User              `json:"updateAuthor"`
	Comment          string             `json:"comment"`
	Created          string             `json:"created"`
	Updated          helper.Iso8601Time `json:"updated"`
	Started          helper.Iso8601Time `json:"started"`
	TimeSpent        string             `json:"timeSpent"`
	TimeSpentSeconds int                `json:"timeSpentSeconds"`
	ID               string             `json:"id"`
	IssueID          uint64             `json:"issueId,string"`
}

func (w Worklog) ToToolLayer(connectionId uint64, issueUpdated *time.Time) *models.JiraWorklog {
	result := &models.JiraWorklog{
		ConnectionId:     connectionId,
		IssueId:          w.IssueID,
		WorklogId:        w.ID,
		TimeSpent:        w.TimeSpent,
		TimeSpentSeconds: w.TimeSpentSeconds,
		Updated:          w.Updated.ToTime(),
		Started:          w.Started.ToTime(),
		IssueUpdated:     issueUpdated,
	}
	if w.Author != nil {
		result.AuthorId = w.Author.EmailAddress
	}
	if w.UpdateAuthor != nil {
		result.UpdateAuthorId = w.UpdateAuthor.EmailAddress
	}
	return result
}
