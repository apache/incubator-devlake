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
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiMrCommitsMeta)
}

var ExtractApiMrCommitsMeta = plugin.SubTaskMeta{
	Name:             "Extract MR Commits",
	EntryPoint:       ExtractApiMergeRequestsCommits,
	EnabledByDefault: true,
	Description:      "Extract raw merge requests commit data into tool layer table GitlabMrCommit and GitlabCommit",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&CollectApiMrCommitsMeta},
}

type GitlabApiCommit struct {
	GitlabId       string `json:"id"`
	Title          string
	Message        string
	ProjectId      int
	ShortId        string             `json:"short_id"`
	AuthorName     string             `json:"author_name"`
	AuthorEmail    string             `json:"author_email"`
	AuthoredDate   common.Iso8601Time `json:"authored_date"`
	CommitterName  string             `json:"committer_name"`
	CommitterEmail string             `json:"committer_email"`
	CommittedDate  common.Iso8601Time `json:"committed_date"`
	WebUrl         string             `json:"web_url"`
	Stats          struct {
		Additions int
		Deletions int
		Total     int
	}
}

func ExtractApiMergeRequestsCommits(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_MERGE_REQUEST_COMMITS_TABLE)

	extractor, err := api.NewStatefulApiExtractor(&api.StatefulApiExtractorArgs[GitlabApiCommit]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Extract: func(gitlabApiCommit *GitlabApiCommit, row *api.RawData) ([]interface{}, errors.Error) {
			// create gitlab commit
			gitlabCommit, err := ConvertCommit(gitlabApiCommit)
			if err != nil {
				return nil, err
			}
			// get input info
			input := &GitlabInput{}
			err = errors.Convert(json.Unmarshal(row.Input, input))
			if err != nil {
				return nil, err
			}

			gitlabMrCommit := &models.GitlabMrCommit{
				CommitSha:          gitlabApiCommit.GitlabId,
				MergeRequestId:     input.GitlabId,
				ConnectionId:       data.Options.ConnectionId,
				CommitAuthorEmail:  gitlabApiCommit.AuthorEmail,
				CommitAuthorName:   gitlabApiCommit.AuthorName,
				CommitAuthoredDate: common.Iso8601TimeToTime(&gitlabApiCommit.AuthoredDate),
			}
			gitlabProjectCommit := &models.GitlabProjectCommit{
				ConnectionId:    data.Options.ConnectionId,
				GitlabProjectId: data.Options.ProjectId,
				CommitSha:       gitlabCommit.Sha,
			}

			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 3)

			results = append(results, gitlabCommit, gitlabProjectCommit, gitlabMrCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

// Convert the API response to our DB model instance
func ConvertCommit(commit *GitlabApiCommit) (*models.GitlabCommit, errors.Error) {
	gitlabCommit := &models.GitlabCommit{
		Sha:            commit.GitlabId,
		Title:          commit.Title,
		Message:        commit.Message,
		ShortId:        commit.ShortId,
		AuthorName:     commit.AuthorName,
		AuthorEmail:    commit.AuthorEmail,
		AuthoredDate:   commit.AuthoredDate.ToTime(),
		CommitterName:  commit.CommitterName,
		CommitterEmail: commit.CommitterEmail,
		CommittedDate:  commit.CommittedDate.ToTime(),
		WebUrl:         commit.WebUrl,
		Additions:      commit.Stats.Additions,
		Deletions:      commit.Stats.Deletions,
		Total:          commit.Stats.Total,
	}
	return gitlabCommit, nil
}
