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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertPullRequestCommentsMeta = core.SubTaskMeta{
	Name:             "convertPullRequestComments",
	EntryPoint:       ConvertPullRequestComments,
	EnabledByDefault: true,
	Description:      "ConvertPullRequestComments data from Github api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}

func ConvertPullRequestComments(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Cursor(
		dal.From(&githubModels.GithubPrComment{}),
		dal.Join("left join _tool_github_pull_requests "+
			"on _tool_github_pull_requests.github_id = _tool_github_pull_request_comments.pull_request_id"),
		dal.Where("repo_id = ? and _tool_github_pull_requests.connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	prCommentIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPrComment{})
	prIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})
	accountIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubAccount{})
	prReviewIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPrReview{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubPrComment{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_COMMENTS_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			githubPullRequestComment := inputRow.(*githubModels.GithubPrComment)
			domainPrComment := &code.PullRequestComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: prCommentIdGen.Generate(data.Options.ConnectionId, githubPullRequestComment.GithubId),
				},
				PullRequestId: prIdGen.Generate(data.Options.ConnectionId, githubPullRequestComment.PullRequestId),
				Body:          githubPullRequestComment.Body,
				UserId:        accountIdGen.Generate(data.Options.ConnectionId, githubPullRequestComment.AuthorUserId),
				CreatedDate:   githubPullRequestComment.GithubCreatedAt,
				CommitSha:     githubPullRequestComment.CommitSha,
				Type:          githubPullRequestComment.Type,
				ReviewId:      prReviewIdGen.Generate(data.Options.ConnectionId, githubPullRequestComment.ReviewId),
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
