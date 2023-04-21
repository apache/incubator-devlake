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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
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
} = (*Gitlab)(nil)

type Gitlab string

func (p Gitlab) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p Gitlab) Connection() interface{} {
	return &models.GitlabConnection{}
}

func (p Gitlab) Scope() interface{} {
	return &models.GitlabProject{}
}

func (p Gitlab) TransformationRule() interface{} {
	return &models.GitlabTransformationRule{}
}

func (p Gitlab) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (p Gitlab) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.GitlabConnection{},
		&models.GitlabAccount{},
		&models.GitlabCommit{},
		&models.GitlabIssue{},
		&models.GitlabIssueLabel{},
		&models.GitlabJob{},
		&models.GitlabMergeRequest{},
		&models.GitlabMrComment{},
		&models.GitlabMrCommit{},
		&models.GitlabMrLabel{},
		&models.GitlabMrNote{},
		&models.GitlabPipeline{},
		&models.GitlabPipelineProject{},
		&models.GitlabProject{},
		&models.GitlabProjectCommit{},
		&models.GitlabReviewer{},
		&models.GitlabTag{},
	}
}

func (p Gitlab) Description() string {
	return "To collect and enrich data from Gitlab"
}

func (p Gitlab) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,
		tasks.CollectApiMergeRequestsMeta,
		tasks.ExtractApiMergeRequestsMeta,
		tasks.CollectApiMergeRequestDetailsMeta,
		tasks.CollectApiMergeRequestDetailsMeta,
		tasks.CollectApiMrNotesMeta,
		tasks.ExtractApiMrNotesMeta,
		tasks.CollectApiMrCommitsMeta,
		tasks.ExtractApiMrCommitsMeta,
		tasks.CollectApiPipelinesMeta,
		tasks.ExtractApiPipelinesMeta,
		tasks.CollectApiPipelineDetailsMeta,
		tasks.ExtractApiPipelineDetailsMeta,
		tasks.CollectApiJobsMeta,
		tasks.ExtractApiJobsMeta,
		tasks.EnrichMergeRequestsMeta,
		tasks.CollectAccountsMeta,
		tasks.ExtractAccountsMeta,
		tasks.ConvertAccountsMeta,
		tasks.ConvertProjectMeta,
		tasks.ConvertApiMergeRequestsMeta,
		tasks.ConvertMrCommentMeta,
		tasks.ConvertApiMrCommitsMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertIssueLabelsMeta,
		tasks.ConvertMrLabelsMeta,
		tasks.ConvertCommitsMeta,
		tasks.ConvertPipelineMeta,
		tasks.ConvertPipelineCommitMeta,
		tasks.ConvertJobMeta,
		tasks.CollectApiCommitsMeta,
		tasks.ExtractApiCommitsMeta,
		tasks.ExtractApiMergeRequestDetailsMeta,
		tasks.CollectTagMeta,
		tasks.ExtractTagMeta,
	}
}

func (p Gitlab) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
	}
	connection := &models.GitlabConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "connection not found")
	}

	apiClient, err := tasks.NewGitlabApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	var timeAfter time.Time
	if op.TimeAfter != "" {
		timeAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.TimeAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `timeAfter`")
		}
	}

	if op.ProjectId != 0 {
		var scope *models.GitlabProject
		// support v100 & advance mode
		// If we still cannot find the record in db, we have to request from remote server and save it to db
		db := taskCtx.GetDal()
		err = db.First(&scope, dal.Where("connection_id = ? AND gitlab_id = ?", op.ConnectionId, op.ProjectId))
		if err != nil && db.IsErrorNotFound(err) {
			var project *models.GitlabApiProject
			project, err = api.GetApiProject(op, apiClient)
			if err != nil {
				return nil, err
			}
			logger.Debug(fmt.Sprintf("Current project: %d", project.GitlabId))
			i := project.ConvertApiScope()
			scope = i.(*models.GitlabProject)
			scope.ConnectionId = op.ConnectionId
			err = taskCtx.GetDal().CreateIfNotExist(&scope)
			if err != nil {
				return nil, err
			}
		}
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find project: %d", op.ProjectId))
		}
	}

	if op.GitlabTransformationRule == nil && op.TransformationRuleId != 0 {
		var transformationRule models.GitlabTransformationRule
		db := taskCtx.GetDal()
		err = db.First(&transformationRule, dal.Where("id = ?", op.TransformationRuleId))
		if err != nil {
			if db.IsErrorNotFound(err) {
				return nil, errors.Default.Wrap(err, fmt.Sprintf("can not find transformationRules by transformationRuleId [%d]", op.TransformationRuleId))
			}
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find transformationRules by transformationRuleId [%d]", op.TransformationRuleId))
		}
		op.GitlabTransformationRule = &transformationRule
	}

	regexEnricher := helper.NewRegexEnricher()
	if err := regexEnricher.TryAdd(devops.DEPLOYMENT, op.DeploymentPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `deploymentPattern`")
	}
	if err := regexEnricher.TryAdd(devops.PRODUCTION, op.ProductionPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `productionPattern`")
	}

	taskData := tasks.GitlabTaskData{
		Options:       op,
		ApiClient:     apiClient,
		RegexEnricher: regexEnricher,
	}

	if !timeAfter.IsZero() {
		taskData.TimeAfter = &timeAfter
		logger.Debug("collect data updated timeAfter %s", timeAfter)
	}
	return &taskData, nil
}

func (p Gitlab) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gitlab"
}

func (p Gitlab) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Gitlab) MakePipelinePlan(connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(p.SubTaskMetas(), connectionId, scope)
}

func (p Gitlab) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
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
		"connections/:connectionId/scopes/:scopeId": {
			"GET":   api.GetScope,
			"PATCH": api.UpdateScope,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": api.SearchRemoteScopes,
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
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
	}
}

func (p Gitlab) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.GitlabTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
