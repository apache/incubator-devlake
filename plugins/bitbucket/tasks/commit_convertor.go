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
	bitbucketModels "github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertCommitsMeta = core.SubTaskMeta{
	Name:             "convertCommits",
	EntryPoint:       ConvertCommits,
	EnabledByDefault: false,
	Required:         false,
	Description:      "Convert tool layer table bitbucket_commits into  domain layer table commits",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}

func ConvertCommits(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)
	repoId := data.Repo.BitbucketId

	cursor, err := db.Cursor(
		dal.From("_tool_bitbucket_commits gc"),
		dal.Join(`left join _tool_bitbucket_repo_commits grc on (
			grc.commit_sha = gc.sha
		)`),
		dal.Select("gc.*"),
		dal.Where("grc.repo_id = ? AND grc.connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	repoDidGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketRepo{})
	domainRepoId := repoDidGen.Generate(data.Options.ConnectionId, repoId)

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_COMMIT_TABLE,
		},
		InputRowType: reflect.TypeOf(bitbucketModels.BitbucketCommit{}),
		Input:        cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			bitbucketCommit := inputRow.(*bitbucketModels.BitbucketCommit)
			domainCommit := &code.Commit{
				Sha:           bitbucketCommit.Sha,
				Message:       bitbucketCommit.Message,
				Additions:     bitbucketCommit.Additions,
				Deletions:     bitbucketCommit.Deletions,
				AuthorId:      bitbucketCommit.AuthorEmail,
				AuthorName:    bitbucketCommit.AuthorName,
				AuthorEmail:   bitbucketCommit.AuthorEmail,
				AuthoredDate:  bitbucketCommit.AuthoredDate,
				CommittedDate: bitbucketCommit.CommittedDate,
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
