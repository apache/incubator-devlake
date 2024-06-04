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

package impl

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/issue_trace/api"
	"github.com/apache/incubator-devlake/plugins/issue_trace/models"
	"github.com/apache/incubator-devlake/plugins/issue_trace/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/issue_trace/services"
	"github.com/apache/incubator-devlake/plugins/issue_trace/tasks"
	"github.com/mitchellh/mapstructure"
)

type IssueTrace struct{}

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginModel
	plugin.PluginMigration
	plugin.PluginApi
} = (*IssueTrace)(nil)

func (p IssueTrace) Name() string {
	return "issue_trace"
}

func (p IssueTrace) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p IssueTrace) Description() string {
	return "To enrich data from issue tracking domain"
}

// Register all subtasks
func (p IssueTrace) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		// issue_status_history
		tasks.ConvertIssueStatusHistoryMeta,
		// issue_assignee_history
		tasks.ConvertIssueAssigneeHistoryMeta,
	}
}

// Prepare your apiClient which will be used to request remote api,
// `apiClient` is defined in `client.go` under `tasks`
// `SprintPerformanceEnricherTaskData` is defined in `task_data.go` under `tasks`
func (p IssueTrace) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
	var op tasks.Options
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, errors.Default.Wrap(err, "Failed to decode options")
	}
	var boardId string
	if op.LakeBoardId != "" {
		boardId = op.LakeBoardId
	} else {
		boardModel := services.GetTicketBoardModel(op.Plugin)
		if boardModel == nil {
			err := errors.BadInput.New("unsupported board type")
			logger.Error(err, "")
			return nil, err
		}
		boardIdGen := didgen.NewDomainIdGenerator(boardModel)
		boardId = boardIdGen.Generate(op.ConnectionId, op.BoardId)
	}

	var taskData = &tasks.TaskData{
		Options: op,
		BoardId: boardId,
	}

	taskData.Options = op
	return taskData, nil
}

func (p IssueTrace) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/issue_trace"
}

func (p IssueTrace) MigrationScripts() []plugin.MigrationScript {
	return []plugin.MigrationScript{
		&migrationscripts.NewIssueTable{},
	}
}

func (p IssueTrace) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{}
}

func (p IssueTrace) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.IssueAssigneeHistory{},
		&models.IssueStatusHistory{},
	}
}
