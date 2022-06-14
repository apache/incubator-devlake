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
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertCommitsMeta = core.SubTaskMeta{
	Name:             "convertCommits",
	EntryPoint:       ConvertCommits,
	EnabledByDefault: false,
	Description:      "Convert tool layer table github_commits into  domain layer table commits",
}

func ConvertCommits(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Cursor(
		dal.From("_tool_github_commits gc"),
		dal.Join(`left join _tool_github_repo_commits grc on (
			grc.commit_sha = gc.sha
		)`),
		dal.Select("gc.*"),
		dal.Where("grc.repo_id = ?", repoId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	repoDidGen := didgen.NewDomainIdGenerator(&githubModels.GithubRepo{})
	domainRepoId := repoDidGen.Generate(repoId)
	userDidGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_COMMENTS_TABLE,
		},
		InputRowType: reflect.TypeOf(githubModels.GithubCommit{}),
		Input:        cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			githubCommit := inputRow.(*githubModels.GithubCommit)
			domainCommit := &code.Commit{
				Sha:            githubCommit.Sha,
				Message:        githubCommit.Message,
				Additions:      githubCommit.Additions,
				Deletions:      githubCommit.Deletions,
				AuthorId:       userDidGen.Generate(githubCommit.AuthorId),
				AuthorName:     githubCommit.AuthorName,
				AuthorEmail:    githubCommit.AuthorEmail,
				AuthoredDate:   githubCommit.AuthoredDate,
				CommitterName:  githubCommit.CommitterName,
				CommitterEmail: githubCommit.CommitterEmail,
				CommittedDate:  githubCommit.CommittedDate,
				CommitterId:    userDidGen.Generate(githubCommit.CommitterId),
			}
			repoCommit := &code.RepoCommit{
				RepoId:    domainRepoId,
				CommitSha: domainCommit.Sha,
			}
			return []interface{}{
				domainCommit,
				repoCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
