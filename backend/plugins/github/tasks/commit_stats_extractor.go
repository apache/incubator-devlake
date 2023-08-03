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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiCommitStatsMeta)
}

var ExtractApiCommitStatsMeta = plugin.SubTaskMeta{
	Name:             "extractApiCommitStats",
	EntryPoint:       ExtractApiCommitStats,
	EnabledByDefault: false,
	Description:      "Extract raw commit stats data into tool layer table github_commit_stats",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	DependencyTables: []string{RAW_COMMIT_STATS_TABLE},
	ProductTables: []string{
		models.GithubCommit{}.TableName(),
		models.GithubCommitStat{}.TableName()},
}

type ApiSingleCommitResponse struct {
	Sha   string
	Stats struct {
		Additions int
		Deletions int
	}
	Commit struct {
		Committer struct {
			Name  string
			Email string
			Date  api.Iso8601Time
		}
	}
}

func ExtractApiCommitStats(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraCommits by Board
			*/
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			/*
				Table store raw data
			*/
			Table: RAW_COMMIT_STATS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			body := &ApiSingleCommitResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}
			if body.Sha == "" {
				return nil, nil
			}

			db := taskCtx.GetDal()
			commit := &models.GithubCommit{}
			err = db.First(commit, dal.Where("sha = ?", body.Sha), dal.Limit(1))
			if err != nil && !db.IsErrorNotFound(err) {
				return nil, err
			}

			commit.Additions = body.Stats.Additions
			commit.Deletions = body.Stats.Deletions

			commitStat := &models.GithubCommitStat{
				ConnectionId:  data.Options.ConnectionId,
				Additions:     body.Stats.Additions,
				Deletions:     body.Stats.Deletions,
				CommittedDate: body.Commit.Committer.Date.ToTime(),
				Sha:           body.Sha,
			}

			results := make([]interface{}, 0, 2)

			results = append(results, commit)
			results = append(results, commitStat)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
