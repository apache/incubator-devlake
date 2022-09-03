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
	"github.com/apache/incubator-devlake/errors"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertCommitsMeta = core.SubTaskMeta{
	Name:             "convertApiCommits",
	EntryPoint:       ConvertCommits,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitee_commits into  domain layer table commits",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_CROSS},
}

func ConvertCommits(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)
	db := taskCtx.GetDal()
	repoId := data.Repo.GiteeId

	// select all commits belongs to the project
	cursor, err := db.Cursor(
		dal.Select("gc.*"),
		dal.From("_tool_gitee_commits gc"),
		dal.Join(`left join _tool_gitee_repo_commits grc on (
			grc.commit_sha = gc.sha
		)`),
		dal.Where("grc.repo_id = ? AND grc.connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	accountIdGen := didgen.NewDomainIdGenerator(&models.GiteeAccount{})
	repoDidGen := didgen.NewDomainIdGenerator(&models.GiteeRepo{})
	domainRepoId := repoDidGen.Generate(data.Options.ConnectionId, repoId)

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GiteeCommit{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			giteeCommit := inputRow.(*models.GiteeCommit)

			// convert commit
			commit := &code.Commit{}
			commit.Sha = giteeCommit.Sha
			commit.Message = giteeCommit.Message
			commit.Additions = giteeCommit.Additions
			commit.Deletions = giteeCommit.Deletions
			commit.AuthorId = accountIdGen.Generate(data.Options.ConnectionId, giteeCommit.AuthorId)
			commit.AuthorName = giteeCommit.AuthorName
			commit.AuthorEmail = giteeCommit.AuthorEmail
			commit.AuthoredDate = giteeCommit.AuthoredDate
			commit.CommitterName = giteeCommit.CommitterName
			commit.CommitterEmail = giteeCommit.CommitterEmail
			commit.CommittedDate = giteeCommit.CommittedDate
			commit.CommitterId = accountIdGen.Generate(data.Options.ConnectionId, giteeCommit.CommitterId)

			// convert repo / commits relationship
			repoCommit := &code.RepoCommit{
				RepoId:    domainRepoId,
				CommitSha: giteeCommit.Sha,
			}

			return []interface{}{
				commit,
				repoCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
