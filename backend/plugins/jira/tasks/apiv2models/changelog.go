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
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type Changelog struct {
	ID      uint64             `json:"id,string"`
	Author  Account            `json:"author"`
	Created common.Iso8601Time `json:"created"`
	Items   []ChangelogItem    `json:"items"`
}

func (c Changelog) ToToolLayer(connectionId, issueId uint64, issueUpdated *time.Time) (*models.JiraIssueChangelogs, *models.JiraAccount) {
	changelog := &models.JiraIssueChangelogs{
		ConnectionId:      connectionId,
		ChangelogId:       c.ID,
		IssueId:           issueId,
		AuthorAccountId:   c.Author.getAccountId(),
		AuthorDisplayName: c.Author.DisplayName,
		AuthorActive:      c.Author.Active,
		Created:           c.Created.ToTime(),
		IssueUpdated:      issueUpdated,
	}
	return changelog, c.Author.ToToolLayer(connectionId)
}

type ChangelogItem struct {
	Field     string `json:"field"`
	Fieldtype string `json:"fieldtype"`
	FieldId   string `json:"fieldId"`

	FromValue  string `json:"from"`
	FromString string `json:"fromString"`

	ToValue  string `json:"to"`
	ToString string `json:"toString"`

	TmpFromAccountId string `json:"tmpFromAccountId,omitempty"`
	TmpToAccountId   string `json:"tmpToAccountId,omitempty"`
}

func (c ChangelogItem) ToToolLayer(connectionId, changelogId uint64) *models.JiraIssueChangelogItems {
	item := &models.JiraIssueChangelogItems{
		ConnectionId:     connectionId,
		ChangelogId:      changelogId,
		Field:            c.Field,
		FieldType:        c.Fieldtype,
		FieldId:          c.FieldId,
		FromValue:        c.FromValue,
		FromString:       c.FromString,
		ToValue:          c.ToValue,
		ToString:         c.ToString,
		TmpFromAccountId: c.TmpFromAccountId,
		TmpToAccountId:   c.TmpToAccountId,
	}
	return item
}

func (c ChangelogItem) ExtractUser(connectionId uint64, userFieldMaps map[string]struct{}) []*models.JiraAccount {
	var result []*models.JiraAccount
	_, ok := userFieldMaps[c.Field]
	if c.Field == "assignee" || c.Field == "reporter" || ok {
		if c.FromValue != "" {
			result = append(result, &models.JiraAccount{ConnectionId: connectionId, AccountId: c.FromValue})
		}
		if c.ToValue != "" {
			result = append(result, &models.JiraAccount{ConnectionId: connectionId, AccountId: c.ToValue})
		}
	}
	return result
}
