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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ConvertExecutionStory

var ConvertExecutionStoryMeta = plugin.SubTaskMeta{
	Name:             "convertExecutionStory",
	EntryPoint:       ConvertExecutionStory,
	EnabledByDefault: true,
	Description:      "convert Zentao execution_stories",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertExecutionStory(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()
	executionIdGen := didgen.NewDomainIdGenerator(&models.ZentaoExecution{})
	storyIdGen := didgen.NewDomainIdGenerator(&models.ZentaoStory{})
	cursor, err := db.Cursor(
		dal.From(&models.ZentaoExecutionStory{}),
		dal.Where(`project_id = ? and connection_id = ?`, data.Options.ProjectId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	convertor, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ZentaoExecutionStory{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_EXECUTION_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			executionStory := inputRow.(*models.ZentaoExecutionStory)
			sprintIssue := &ticket.SprintIssue{
				SprintId: executionIdGen.Generate(data.Options.ConnectionId, executionStory.ExecutionId),
				IssueId:  storyIdGen.Generate(data.Options.ConnectionId, executionStory.StoryId),
			}
			return []interface{}{sprintIssue}, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
