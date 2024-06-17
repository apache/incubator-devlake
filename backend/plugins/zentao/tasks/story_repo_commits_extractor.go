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
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ExtractStoryRepoCommits

var ExtractStoryRepoCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractStoryRepoCommits",
	EntryPoint:       ExtractStoryRepoCommits,
	EnabledByDefault: false,
	Description:      "extract Zentao story repo commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractStoryRepoCommits(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_STORY_REPO_COMMITS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			res := &models.ZentaoStoryRepoCommitsRes{}
			err := json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}

			results := make([]interface{}, 0)
			issueIds, err := extractIdFromLogComment("story", res.Log.Comment)
			if err != nil {
				return nil, errors.Default.Wrap(err, "extractIdFromLogComment")
			}
			for _, issueId := range issueIds {
				storyRepoCommits := &models.ZentaoStoryRepoCommit{
					ConnectionId: data.Options.ConnectionId,
					Project:      data.Options.ProjectId,
					RepoUrl:      res.Repo.CodePath,
					CommitSha:    res.Revision,
					IssueId:      issueId,
				}
				results = append(results, storyRepoCommits)
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
