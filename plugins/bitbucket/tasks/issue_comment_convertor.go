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
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	bitbucketModels "github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertIssueCommentsMeta = core.SubTaskMeta{
	Name:             "convertIssueComments",
	EntryPoint:       ConvertIssueComments,
	EnabledByDefault: true,
	Description:      "ConvertIssueComments data from Bitbucket api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}

func ConvertIssueComments(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)
	repoId := data.Repo.BitbucketId

	cursor, err := db.Cursor(
		dal.From(&bitbucketModels.BitbucketIssueComment{}),
		dal.Join("left join _tool_bitbucket_issues "+
			"on _tool_bitbucket_issues.bitbucket_id = _tool_bitbucket_issue_comments.issue_id"),
		dal.Where("repo_id = ? and _tool_bitbucket_issues.connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketIssue{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(bitbucketModels.BitbucketIssueComment{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_ISSUE_COMMENTS_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			bitbucketIssueComment := inputRow.(*bitbucketModels.BitbucketIssueComment)
			domainIssueComment := &ticket.IssueComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: issueIdGen.Generate(data.Options.ConnectionId, bitbucketIssueComment.BitbucketId),
				},
				IssueId:     issueIdGen.Generate(data.Options.ConnectionId, bitbucketIssueComment.IssueId),
				UserId:      bitbucketIssueComment.AuthorUserId,
				CreatedDate: bitbucketIssueComment.CreatedAt,
			}
			return []interface{}{
				domainIssueComment,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
