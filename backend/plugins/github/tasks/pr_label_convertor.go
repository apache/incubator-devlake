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
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertPullRequestLabelsMeta)
}

var ConvertPullRequestLabelsMeta = plugin.SubTaskMeta{
	Name:             "convertPullRequestLabels",
	EntryPoint:       ConvertPullRequestLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_pull_request_labels into  domain layer table pull_request_labels",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{
		models.GithubPrLabel{}.TableName(),     // cursor
		models.GithubPullRequest{}.TableName(), // cursor and id generator
		RAW_PULL_REQUEST_TABLE},
	ProductTables: []string{code.PullRequestLabel{}.TableName()},
}

func ConvertPullRequestLabels(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Options.GithubId

	cursor, err := db.Cursor(
		dal.From(&models.GithubPrLabel{}),
		dal.Join(`left join _tool_github_pull_requests on _tool_github_pull_requests.github_id = _tool_github_pull_request_labels.pull_id`),
		dal.Where("_tool_github_pull_requests.repo_id = ? and _tool_github_pull_requests.connection_id = ?", repoId, data.Options.ConnectionId),
		dal.Orderby("pull_id ASC"),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	prIdGen := didgen.NewDomainIdGenerator(&models.GithubPullRequest{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.GithubPrLabel{}),
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
			prLabel := inputRow.(*models.GithubPrLabel)
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
