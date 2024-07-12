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
	"encoding/json"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

var _ plugin.SubTaskEntryPoint = ExtractIssueFields

var ExtractIssueFieldsMeta = plugin.SubTaskMeta{
	Name:             "extractIssueFields",
	EntryPoint:       ExtractIssueFields,
	EnabledByDefault: true,
	Description:      "extract Jira issue fields",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

type JiraIssueField struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Custom      bool     `json:"custom"`
	Orderable   bool     `json:"orderable"`
	Navigable   bool     `json:"navigable"`
	Searchable  bool     `json:"searchable"`
	ClauseNames []string `json:"clauseNames"`
	Schema      struct {
		Type     string `json:"type"`
		Items    string `json:"items"`
		Custom   string `json:"custom"`
		System   string `json:"system"`
		CustomID int    `json:"customId"`
	} `json:"schema"`
}

func ExtractIssueFields(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_ISSUE_FIELDS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var issueField JiraIssueField
			err := errors.Convert(json.Unmarshal(row.Data, &issueField))
			if err != nil {
				return nil, err
			}
			jiraIssueField := &models.JiraIssueField{
				NoPKModel:    common.NewNoPKModel(),
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,

				ID:         issueField.ID,
				Name:       issueField.Name,
				Custom:     issueField.Custom,
				Orderable:  issueField.Orderable,
				Navigable:  issueField.Navigable,
				Searchable: issueField.Searchable,
				//ClauseNames:      issueField.ClauseNames,
				SchemaType:       issueField.Schema.Type,
				SchemaItems:      issueField.Schema.Items,
				SchemaCustom:     issueField.Schema.Custom,
				SchemaCustomID:   issueField.Schema.CustomID,
				ScheCustomSystem: issueField.Schema.System,
			}
			return []interface{}{jiraIssueField}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
