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
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

func ConvertTaskCommit(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_COMMIT_TABLE, false)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("convert board:%d", data.Options.WorkspaceId)

	clauses := []dal.Clause{
		dal.From(&models.TapdTaskCommit{}),
		dal.Where("connection_id = ? AND workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	issueIdGen := didgen.NewDomainIdGenerator(&models.TapdIssue{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TapdTaskCommit{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolL := inputRow.(*models.TapdTaskCommit)
			domainL := &crossdomain.IssueCommit{
				IssueId:   issueIdGen.Generate(data.Options.ConnectionId, toolL.TaskId),
				CommitSha: toolL.CommitId,
			}

			return []interface{}{
				domainL,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertTaskCommitMeta = core.SubTaskMeta{
	Name:             "convertTaskCommit",
	EntryPoint:       ConvertTaskCommit,
	EnabledByDefault: true,
	Description:      "convert Tapd TaskCommit",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}
