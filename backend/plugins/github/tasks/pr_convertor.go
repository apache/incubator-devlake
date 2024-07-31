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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertPullRequestsMeta)
}

var ConvertPullRequestsMeta = plugin.SubTaskMeta{
	Name:             "Convert Pull Requests",
	EntryPoint:       ConvertPullRequests,
	EnabledByDefault: true,
	Description:      "ConvertPullRequests data from Github api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS, plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{
		models.GithubPullRequest{}.TableName(), // cursor
		//models.GithubRepo{}.TableName(),        // id generator, but not regard as dependency
		models.GithubAccount{}.TableName(), // cursor
		RAW_PULL_REQUEST_TABLE},
	ProductTables: []string{code.PullRequest{}.TableName()},
}

func ConvertPullRequests(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Options.GithubId

	cursor, err := db.Cursor(
		dal.From(&models.GithubPullRequest{}),
		dal.Where("repo_id = ? and connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	prIdGen := didgen.NewDomainIdGenerator(&models.GithubPullRequest{})
	repoIdGen := didgen.NewDomainIdGenerator(&models.GithubRepo{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.GithubAccount{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.GithubPullRequest{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_PULL_REQUEST_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			pr := inputRow.(*models.GithubPullRequest)
			domainPr := &code.PullRequest{
				DomainEntity: domainlayer.DomainEntity{
					Id: prIdGen.Generate(data.Options.ConnectionId, pr.GithubId),
				},
				BaseRepoId:     repoIdGen.Generate(data.Options.ConnectionId, pr.RepoId),
				HeadRepoId:     repoIdGen.Generate(data.Options.ConnectionId, pr.HeadRepoId),
				OriginalStatus: pr.State,
				Title:          pr.Title,
				Url:            pr.Url,
				AuthorName:     pr.AuthorName,
				Description:    pr.Body,
				CreatedDate:    pr.GithubCreatedAt,
				MergedDate:     pr.MergedAt,
				ClosedDate:     pr.ClosedAt,
				PullRequestKey: pr.Number,
				Type:           pr.Type,
				Component:      pr.Component,
				MergeCommitSha: pr.MergeCommitSha,
				BaseRef:        pr.BaseRef,
				BaseCommitSha:  pr.BaseCommitSha,
				HeadRef:        pr.HeadRef,
				HeadCommitSha:  pr.HeadCommitSha,
				Additions:      pr.Additions,
				Deletions:      pr.Deletions,
				MergedByName:   pr.MergedByName,
				MergedById:     accountIdGen.Generate(data.Options.ConnectionId, pr.MergedById),
				IsDraft:        pr.IsDraft,
			}
			if pr.State == "open" || pr.State == "OPEN" {
				domainPr.Status = code.OPEN
			} else if (pr.State == "closed" && pr.MergedAt != nil) || pr.State == "MERGED" {
				domainPr.Status = code.MERGED
			} else {
				domainPr.Status = code.CLOSED
			}
			if pr.AuthorId != 0 {
				domainPr.AuthorId = accountIdGen.Generate(data.Options.ConnectionId, pr.AuthorId)
			}
			return []interface{}{
				domainPr,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
