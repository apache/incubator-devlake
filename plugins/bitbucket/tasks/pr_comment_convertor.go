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

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	bitbucketModels "github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertPrCommentsMeta = core.SubTaskMeta{
	Name:             "convertPullRequestComments",
	EntryPoint:       ConvertPullRequestComments,
	EnabledByDefault: true,
	Description:      "ConvertPullRequestComments data from Bitbucket api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE_REVIEW},
}

func ConvertPullRequestComments(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)
	repoId := data.Repo.BitbucketId

	cursor, err := db.Cursor(
		dal.From(&bitbucketModels.BitbucketPrComment{}),
		dal.Join("left join _tool_bitbucket_pull_requests "+
			"on _tool_bitbucket_pull_requests.bitbucket_id = _tool_bitbucket_pull_request_comments.pull_request_id"),
		dal.Where("repo_id = ? and _tool_bitbucket_pull_requests.connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	domainIdGeneratorComment := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketPrComment{})
	prIdGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketPullRequest{})
	accountIdGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(bitbucketModels.BitbucketPrComment{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_PULL_REQUEST_COMMENTS_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			bitbucketPullRequestComment := inputRow.(*bitbucketModels.BitbucketPrComment)
			domainPrComment := &code.PullRequestComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: domainIdGeneratorComment.Generate(data.Options.ConnectionId, bitbucketPullRequestComment.BitbucketId),
				},
				PullRequestId: prIdGen.Generate(data.Options.ConnectionId, bitbucketPullRequestComment.PullRequestId),
				UserId:        accountIdGen.Generate(data.Options.ConnectionId, bitbucketPullRequestComment.AuthorId),
				CreatedDate:   bitbucketPullRequestComment.CreatedAt,
				Body:          bitbucketPullRequestComment.Body,
				Type:          bitbucketPullRequestComment.Type,
				CommitSha:     "",
				Position:      0,
			}
			return []interface{}{
				domainPrComment,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
