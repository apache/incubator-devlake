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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"time"
)

type Comment struct {
	Self         string             `json:"self"`
	Id           string             `json:"id"`
	Author       *Account           `json:"author"`
	Body         string             `json:"body"`
	UpdateAuthor *Account           `json:"updateAuthor"`
	Created      helper.Iso8601Time `json:"created"`
	Updated      helper.Iso8601Time `json:"updated"`
	JsdPublic    bool               `json:"jsdPublic"`
}

func (c Comment) ToToolLayer(connectionId uint64, issueId uint64, issueUpdated *time.Time) *models.JiraIssueComment {
	result := &models.JiraIssueComment{
		ConnectionId: connectionId,
		IssueId:      issueId,
		ComentId:     c.Id,
		Self:         c.Self,
		Body:         c.Body,
		Created:      c.Updated.ToTime(),
		Updated:      c.Updated.ToTime(),
		IssueUpdated: issueUpdated,
	}
	if c.Author != nil {
		result.CreatorAccountId = c.Author.getAccountId()
		result.CreatorDisplayName = c.Author.DisplayName
	}
	return result
}
