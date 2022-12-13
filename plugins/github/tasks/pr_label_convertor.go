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
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertPullRequestLabelsMeta = core.SubTaskMeta{
	Name:             "convertPullRequestLabels",
	EntryPoint:       ConvertPullRequestLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_pull_request_labels into  domain layer table pull_request_labels",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE_REVIEW},
}

func ConvertPullRequestLabels(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Options.GithubId

	cursor, err := db.Cursor(
		dal.From(&githubModels.GithubPrLabel{}),
		dal.Join(`left join _tool_github_pull_requests on _tool_github_pull_requests.github_id = _tool_github_pull_request_labels.pull_id`),
		dal.Where("_tool_github_pull_requests.repo_id = ? and _tool_github_pull_requests.connection_id = ?", repoId, data.Options.ConnectionId),
		dal.Orderby("pull_id ASC"),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	prIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubPrLabel{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_PULL_REQUEST_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			prLabel := inputRow.(*githubModels.GithubPrLabel)
			domainPrLabel := &code.PullRequestLabel{
				PullRequestId: prIdGen.Generate(data.Options.ConnectionId, prLabel.PullId),
				LabelName:     prLabel.LabelName,
			}
			return []interface{}{
				domainPrLabel,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
