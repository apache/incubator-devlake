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
	"regexp"
)

var ExtractRemotelinksMeta = plugin.SubTaskMeta{
	Name:             "extractRemotelinks",
	EntryPoint:       ExtractRemotelinks,
	EnabledByDefault: true,
	Description:      "extract Jira remote links",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractRemotelinks(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("extract remote links")

	var err errors.Error
	var commitShaRegex *regexp.Regexp
	if data.Options.TransformationRules != nil && data.Options.TransformationRules.RemotelinkCommitShaPattern != "" {
		pattern := data.Options.TransformationRules.RemotelinkCommitShaPattern
		commitShaRegex, err = errors.Convert01(regexp.Compile(pattern))
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile pattern failed")
		}
	}
	var commitRepoUrlRegexps []*regexp.Regexp
	if tr := data.Options.TransformationRules; tr != nil {
		for _, s := range tr.RemotelinkRepoPattern {
			pattern, e := regexp.Compile(s)
			if e != nil {
				return errors.Convert(e)
			}
			commitRepoUrlRegexps = append(commitRepoUrlRegexps, pattern)
		}
	}

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_REMOTELINK_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var result []interface{}
			var raw apiv2models.RemoteLink
			err := errors.Convert(json.Unmarshal(row.Data, &raw))
			if err != nil {
				return nil, err
			}
			var input apiv2models.Input
			err = errors.Convert(json.Unmarshal(row.Input, &input))
			if err != nil {
				return nil, err
			}
			remotelink := &models.JiraRemotelink{
				ConnectionId: connectionId,
				RemotelinkId: raw.ID,
				IssueId:      input.IssueId,
				Self:         raw.Self,
				Title:        raw.Object.Title,
				Url:          raw.Object.URL,
				IssueUpdated: &input.UpdateTime,
			}
			result = append(result, remotelink)
			issueCommit := &models.JiraIssueCommit{
				ConnectionId: connectionId,
				IssueId:      remotelink.IssueId,
				CommitUrl:    remotelink.Url,
			}
			if commitShaRegex != nil {
				groups := commitShaRegex.FindStringSubmatch(remotelink.Url)
				if len(groups) > 1 {
					issueCommit.CommitSha = groups[1]
				}
			}
			if issueCommit.CommitSha == "" {
				issueCommit.CommitSha = api.ExtractCommitSha(commitRepoUrlRegexps, remotelink.Url)
			}
			if issueCommit.CommitSha != "" {
				result = append(result, issueCommit)
			}
			return result, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}
