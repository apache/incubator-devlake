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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"reflect"
)

var ConvertCommitsMeta = plugin.SubTaskMeta{
	Name:             "convertCommits",
	EntryPoint:       ConvertCommits,
	EnabledByDefault: false,
	Required:         false,
	Description:      "Convert tool layer table bitbucket_commits into  domain layer table commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

func ConvertCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)
	db := taskCtx.GetDal()
	repoId := data.Options.FullName

	cursor, err := db.Cursor(
		dal.From("_tool_bitbucket_commits c"),
		dal.Join(`left join _tool_bitbucket_repo_commits rc on (
			rc.commit_sha = c.sha
		)`),
		dal.Select("c.*"),
		dal.Where("rc.repo_id = ? AND rc.connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	repoDidGen := didgen.NewDomainIdGenerator(&models.BitbucketRepo{})
	domainRepoId := repoDidGen.Generate(data.Options.ConnectionId, repoId)

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.BitbucketCommit{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bitbucketCommit := inputRow.(*models.BitbucketCommit)
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
