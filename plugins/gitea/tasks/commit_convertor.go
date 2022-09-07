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
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitea/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertCommitsMeta = core.SubTaskMeta{
	Name:             "convertApiCommits",
	EntryPoint:       ConvertCommits,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitea_commits into  domain layer table commits",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_CROSS},
}

func ConvertCommits(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)
	db := taskCtx.GetDal()
	repoId := data.Repo.GiteaId

	// select all commits belongs to the project
	cursor, err := db.Cursor(
		dal.Select("gc.*"),
		dal.From("_tool_gitea_commits gc"),
		dal.Join(`left join _tool_gitea_repo_commits grc on (
			grc.commit_sha = gc.sha
		)`),
		dal.Where("grc.repo_id = ? AND grc.connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	accountIdGen := didgen.NewDomainIdGenerator(&models.GiteaAccount{})
	repoDidGen := didgen.NewDomainIdGenerator(&models.GiteaRepo{})
	domainRepoId := repoDidGen.Generate(data.Options.ConnectionId, repoId)

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GiteaCommit{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			giteaCommit := inputRow.(*models.GiteaCommit)

			// convert commit
			commit := &code.Commit{}
			commit.Sha = giteaCommit.Sha
			commit.Message = giteaCommit.Message
			commit.Additions = giteaCommit.Additions
			commit.Deletions = giteaCommit.Deletions
			commit.AuthorId = accountIdGen.Generate(data.Options.ConnectionId, giteaCommit.AuthorId)
			commit.AuthorName = giteaCommit.AuthorName
			commit.AuthorEmail = giteaCommit.AuthorEmail
			commit.AuthoredDate = giteaCommit.AuthoredDate
			commit.CommitterName = giteaCommit.CommitterName
			commit.CommitterEmail = giteaCommit.CommitterEmail
			commit.CommittedDate = giteaCommit.CommittedDate
			commit.CommitterId = accountIdGen.Generate(data.Options.ConnectionId, giteaCommit.CommitterId)

			// convert repo / commits relationship
			repoCommit := &code.RepoCommit{
				RepoId:    domainRepoId,
				CommitSha: giteaCommit.Sha,
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
