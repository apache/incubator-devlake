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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"reflect"
)

func init() {
	RegisterSubtaskMeta(&ConvertPrLabelsMeta)
}

var ConvertPrLabelsMeta = plugin.SubTaskMeta{
	Name:             "convertPrLabels",
	EntryPoint:       ConvertPrLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table azuredevops_go_pull_request_labels into domain layer table pull_request_labels",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{
		models.AzuredevopsPrLabel{}.TableName(),
	},
}

func ConvertPrLabels(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawPullRequestTable)
	repoId := data.Options.RepositoryId
	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From(&models.AzuredevopsPrLabel{}),
		dal.Join(`left join _tool_azuredevops_go_pull_requests on
			_tool_azuredevops_go_pull_requests.azuredevops_id = _tool_azuredevops_go_pull_request_labels.pull_request_id`),
		dal.Where(`_tool_azuredevops_go_pull_requests.repository_id = ?
			and _tool_azuredevops_go_pull_requests.connection_id = ?`,
			repoId, data.Options.ConnectionId),
		dal.Orderby("pull_request_id ASC"),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	prIdGen := didgen.NewDomainIdGenerator(&models.AzuredevopsPullRequest{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.AzuredevopsPrLabel{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			prLabel := inputRow.(*models.AzuredevopsPrLabel)
			domainIssueLabel := &code.PullRequestLabel{
				PullRequestId: prIdGen.Generate(data.Options.ConnectionId, prLabel.PullRequestId),
				LabelName:     prLabel.LabelName,
			}
			return []interface{}{
				domainIssueLabel,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
