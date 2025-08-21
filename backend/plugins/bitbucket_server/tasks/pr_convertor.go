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
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

var ConvertPullRequestsMeta = plugin.SubTaskMeta{
	Name:             "convertPullRequests",
	EntryPoint:       ConvertPullRequests,
	EnabledByDefault: true,
	Required:         false,
	Description:      "ConvertPullRequests data from Bitbucket api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

func ConvertPullRequests(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_TABLE)
	db := taskCtx.GetDal()
	repoId := data.Options.FullName

	cursor, err := db.Cursor(
		dal.From(&models.BitbucketServerPullRequest{}),
		dal.Where("repo_id = ? and connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	prIdGen := didgen.NewDomainIdGenerator(&models.BitbucketServerPullRequest{})
	repoIdGen := didgen.NewDomainIdGenerator(&models.BitbucketServerRepo{})
	domainUserIdGen := didgen.NewDomainIdGenerator(&models.BitbucketServerUser{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.BitbucketServerPullRequest{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			pr := inputRow.(*models.BitbucketServerPullRequest)
			domainPr := &code.PullRequest{
				DomainEntity: domainlayer.DomainEntity{
					Id: prIdGen.Generate(data.Options.ConnectionId, pr.RepoId, pr.BitbucketId),
				},
				BaseRepoId:     repoIdGen.Generate(data.Options.ConnectionId, pr.BaseRepoId),
				HeadRepoId:     repoIdGen.Generate(data.Options.ConnectionId, pr.HeadRepoId),
				OriginalStatus: pr.State,
				Title:          pr.Title,
				Url:            pr.Url,
				AuthorId:       domainUserIdGen.Generate(data.Options.ConnectionId, pr.AuthorId),
				AuthorName:     pr.AuthorName,
				Description:    pr.Description,
				CreatedDate:    pr.BitbucketServerCreatedAt,
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
			}

			switch pr.State {
			case "OPEN":
				domainPr.Status = code.OPEN
			case "MERGED":
				domainPr.Status = code.MERGED
			case "DECLINED":
				domainPr.Status = code.CLOSED
			default:
				domainPr.Status = pr.State
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
