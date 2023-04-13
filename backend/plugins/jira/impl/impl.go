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
	"fmt"
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginModel
	plugin.PluginMigration
	plugin.PluginBlueprintV100
	plugin.DataSourcePluginBlueprintV200
	plugin.CloseablePluginTask
	plugin.PluginSource
} = (*Jira)(nil)

type Jira struct {
}

func (p Jira) Connection() interface{} {
	return &models.JiraConnection{}
}

func (p Jira) Scope() interface{} {
	return &models.JiraBoard{}
}

func (p Jira) TransformationRule() interface{} {
	return &models.JiraTransformationRule{}
}

func (p *Jira) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p Jira) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.ApiMyselfResponse{},
		&models.JiraAccount{},
		&models.JiraBoard{},
		&models.JiraBoardIssue{},
		&models.JiraBoardSprint{},
		&models.JiraConnection{},
		&models.JiraIssue{},
		&models.JiraIssueChangelogItems{},
		&models.JiraIssueChangelogs{},
		&models.JiraIssueCommit{},
		&models.JiraIssueLabel{},
		&models.JiraIssueType{},
		&models.JiraProject{},
		&models.JiraRemotelink{},
		&models.JiraServerInfo{},
		&models.JiraSprint{},
		&models.JiraSprintIssue{},
		&models.JiraStatus{},
		&models.JiraWorklog{},
	}
}

func (p Jira) Description() string {
	return "To collect and enrich data from JIRA"
}

func (p Jira) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectStatusMeta,
		tasks.ExtractStatusMeta,

		tasks.CollectProjectsMeta,
		tasks.ExtractProjectsMeta,

		tasks.CollectIssueTypesMeta,
		tasks.ExtractIssueTypesMeta,

		tasks.CollectIssuesMeta,
		tasks.ExtractIssuesMeta,

		tasks.ConvertIssueLabelsMeta,

		tasks.CollectIssueChangelogsMeta,
		tasks.ExtractIssueChangelogsMeta,

		tasks.CollectAccountsMeta,

		tasks.CollectWorklogsMeta,
		tasks.ExtractWorklogsMeta,

		tasks.CollectRemotelinksMeta,
		tasks.ExtractRemotelinksMeta,

		tasks.CollectSprintsMeta,
		tasks.ExtractSprintsMeta,

		tasks.ConvertBoardMeta,

		tasks.ConvertIssuesMeta,

		tasks.ConvertWorklogsMeta,

		tasks.ConvertIssueChangelogsMeta,

		tasks.ConvertSprintsMeta,
		tasks.ConvertSprintIssuesMeta,

		tasks.ConvertIssueCommitsMeta,
		tasks.ConvertIssueRepoCommitsMeta,

		tasks.ExtractAccountsMeta,
		tasks.ConvertAccountsMeta,

		tasks.CollectEpicsMeta,
		tasks.ExtractEpicsMeta,
	}
}

func (p Jira) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.JiraOptions
	var err errors.Error
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	err = helper.Decode(options, &op, nil)
	if err != nil {
		return nil, errors.Default.Wrap(err, "could not decode Jira options")
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("jira connectionId is invalid")
	}
	connection := &models.JiraConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Jira connection")
	}
	jiraApiClient, err := tasks.NewJiraApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to create jira api client")
	}

	if op.BoardId != 0 {
		var scope *models.JiraBoard
		// support v100 & advance mode
		// If we still cannot find the record in db, we have to request from remote server and save it to db
		db := taskCtx.GetDal()
		err = db.First(&scope, dal.Where("connection_id = ? AND board_id = ?", op.ConnectionId, op.BoardId))
		if err != nil && db.IsErrorNotFound(err) {
			var board *apiv2models.Board
			board, err = api.GetApiJira(&op, jiraApiClient)
			if err != nil {
				return nil, err
			}
			logger.Debug(fmt.Sprintf("Current project: %d", board.ID))
			scope = board.ToToolLayer(connection.ID)
			err = db.CreateIfNotExist(&scope)
			if err != nil {
				return nil, err
			}
		}
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find board: %d", op.BoardId))
		}
	}

	if op.BoardId == 0 && op.ScopeId != "" {
		var jiraBoard models.JiraBoard
		// get repo from db
		err = db.First(&jiraBoard, dal.Where(`connection_id = ? and board_id = ?`, connection.ID, op.ScopeId))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find board%s", op.ScopeId))
		}
		op.BoardId = jiraBoard.BoardId
		if op.TransformationRuleId == 0 {
			op.TransformationRuleId = jiraBoard.TransformationRuleId
		}
	}
	if op.TransformationRules == nil && op.TransformationRuleId != 0 {
		var transformationRule models.JiraTransformationRule
		err = taskCtx.GetDal().First(&transformationRule, dal.Where("id = ?", op.TransformationRuleId))
		if err != nil && db.IsErrorNotFound(err) {
			return nil, errors.BadInput.Wrap(err, "fail to get transformationRule")
		}
		op.TransformationRules, err = tasks.MakeTransformationRules(transformationRule)
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "fail to make transformationRule")
		}
	}

	// set default page size
	if op.PageSize <= 0 || op.PageSize > 100 {
		op.PageSize = 100
	}

	info, code, err := tasks.GetJiraServerInfo(jiraApiClient)
	if err != nil || code != http.StatusOK || info == nil {
		return nil, errors.HttpStatus(code).Wrap(err, "fail to get Jira server info")
	}
	taskData := &tasks.JiraTaskData{
		Options:        &op,
		ApiClient:      jiraApiClient,
		JiraServerInfo: *info,
	}
	if op.TimeAfter != "" {
		var timeAfter time.Time
		timeAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.TimeAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `timeAfter`")
		}
		taskData.TimeAfter = &timeAfter
		logger.Debug("collect data created from %s", timeAfter)
	}
	return taskData, nil
}

func (p Jira) MakePipelinePlan(connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	return api.MakePipelinePlanV100(p.SubTaskMetas(), connectionId, scope)
}

func (p Jira) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (pp plugin.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (p Jira) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/jira"
}

func (p Jira) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Jira) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"echo": {
			"POST": func(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
				return &plugin.ApiResourceOutput{Body: input.Body}, nil
			},
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
			"GET":    api.GetConnection,
		},
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":   api.GetScope,
			"PATCH": api.UpdateScope,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScope,
		},
		"connections/:connectionId/transformation_rules": {
			"POST": api.CreateTransformationRule,
			"GET":  api.GetTransformationRuleList,
		},
		"connections/:connectionId/transformation_rules/:id": {
			"PATCH": api.UpdateTransformationRule,
			"GET":   api.GetTransformationRule,
		},
	}
}

func (p Jira) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.JiraTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
