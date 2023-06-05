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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ plugin.SubTaskEntryPoint = ExtractIssues

var ExtractDevelopmentPanelMeta = plugin.SubTaskMeta{
	Name:             "ExtractDevelopmentPanel",
	EntryPoint:       ExtractDevelopmentPanel,
	EnabledByDefault: true,
	Description:      "Extract Jira development panel",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET, plugin.DOMAIN_TYPE_CROSS},
}

func ExtractDevelopmentPanel(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	scopeConfig := data.Options.ScopeConfig
	// if the condition is true, it means that the task is not enabled
	if scopeConfig == nil || scopeConfig.ApplicationType == "" {
		return nil
	}
	connectionId := data.Options.ConnectionId
	var err errors.Error
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_DEVELOPMENT_PANEL,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var result []interface{}
			var raw apiv2models.DevelopmentPanel
			err := errors.Convert(json.Unmarshal(row.Data, &raw))
			if err != nil {
				return nil, err
			}
			var input apiv2models.Input
			err = errors.Convert(json.Unmarshal(row.Input, &input))
			if err != nil {
				return nil, err
			}
			for _, item := range raw.Detail {
				for _, repo := range item.Repositories {
					for _, commit := range repo.Commits {
						issueCommit := &models.JiraIssueCommit{
							ConnectionId: connectionId,
							IssueId:      input.IssueId,
							CommitUrl:    commit.URL,
							CommitSha:    commit.ID,
							RepoUrl:      repo.URL,
						}
						if issueCommit.CommitSha != "" {
							result = append(result, issueCommit)
						}
					}
				}
			}

			return result, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
