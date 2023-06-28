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
	"regexp"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ExtractTaskRepoCommits

var ExtractTaskRepoCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractTaskRepoCommits",
	EntryPoint:       ExtractTaskRepoCommits,
	EnabledByDefault: true,
	Description:      "extract Zentao task repo commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractTaskRepoCommits(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)

	// this Extract only work for project
	if data.Options.ProjectId == 0 {
		return nil
	}

	re := regexp.MustCompile(`(\d+)(?:,\s*(\d+))*`)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ScopeParams(
				data.Options.ConnectionId,
				data.Options.ProjectId,
				data.Options.ProductId,
			),
			Table: RAW_TASK_REPO_COMMITS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			res := &models.ZentaoTaskRepoCommitsRes{}
			err := json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}

			results := make([]interface{}, 0)
			match := re.FindStringSubmatch(res.Log.Comment)
			for i := 1; i < len(match); i++ {
				if match[i] != "" {
					taskRepoCommits := &models.ZentaoTaskRepoCommit{
						ConnectionId: data.Options.ConnectionId,
						Product:      data.Options.ProductId,
						Project:      data.Options.ProjectId,
						RepoUrl:      res.Repo.CodePath,
						CommitSha:    res.Revision,
						IssueId:      match[i], // task id
					}
					results = append(results, taskRepoCommits)
				}
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
