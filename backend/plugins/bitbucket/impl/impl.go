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
	"github.com/apache/incubator-devlake/plugins/bitbucket/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
)

var _ plugin.PluginMeta = (*Bitbucket)(nil)
var _ plugin.PluginInit = (*Bitbucket)(nil)
var _ plugin.PluginTask = (*Bitbucket)(nil)
var _ plugin.PluginApi = (*Bitbucket)(nil)
var _ plugin.PluginModel = (*Bitbucket)(nil)
var _ plugin.PluginMigration = (*Bitbucket)(nil)
var _ plugin.CloseablePluginTask = (*Bitbucket)(nil)
var _ plugin.PluginSource = (*Bitbucket)(nil)

type Bitbucket string

func (p Bitbucket) Connection() interface{} {
	return &models.BitbucketConnection{}
}

func (p Bitbucket) Scope() interface{} {
	return &models.BitbucketRepo{}
}

func (p Bitbucket) TransformationRule() interface{} {
	return &models.BitbucketTransformationRule{}
}

func (p Bitbucket) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p Bitbucket) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.BitbucketConnection{},
		&models.BitbucketAccount{},
		&models.BitbucketCommit{},
		&models.BitbucketPullRequest{},
		&models.BitbucketIssue{},
		&models.BitbucketPrComment{},
		&models.BitbucketIssueComment{},
		&models.BitbucketPipeline{},
		&models.BitbucketRepo{},
		&models.BitbucketRepoCommit{},
	}
}

func (p Bitbucket) Description() string {
	return "To collect and enrich data from Bitbucket"
}

func (p Bitbucket) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectApiPullRequestsMeta,
		tasks.ExtractApiPullRequestsMeta,

		tasks.CollectApiPrCommentsMeta,
		tasks.ExtractApiPrCommentsMeta,

		tasks.CollectApiPrCommitsMeta,
		tasks.ExtractApiPrCommitsMeta,

		tasks.CollectApiCommitsMeta,
		tasks.ExtractApiCommitsMeta,

		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,

		tasks.CollectApiIssueCommentsMeta,
		tasks.ExtractApiIssueCommentsMeta,

		tasks.CollectApiPipelinesMeta,
		tasks.ExtractApiPipelinesMeta,

		tasks.CollectApiDeploymentsMeta,
		tasks.ExtractApiDeploymentsMeta,

		// must run after deployment to match
		tasks.CollectPipelineStepsMeta,
		tasks.ExtractPipelineStepsMeta,

		tasks.ConvertRepoMeta,
		tasks.ConvertAccountsMeta,
		tasks.ConvertPullRequestsMeta,
		tasks.ConvertPrCommentsMeta,
		tasks.ConvertPrCommitsMeta,
		tasks.ConvertCommitsMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertIssueCommentsMeta,
		tasks.ConvertPipelineMeta,
		tasks.ConvertPipelineStepMeta,
		tasks.ConvertiDeploymentMeta,
	}
}

func (p Bitbucket) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.BitbucketConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get bitbucket connection by the given connection ID")
	}

	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get bitbucket API client instance")
	}
	err = EnrichOptions(taskCtx, op, apiClient.ApiClient)
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
	regexEnricher := helper.NewRegexEnricher()
	if err := regexEnricher.TryAdd(devops.DEPLOYMENT, op.DeploymentPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `deploymentPattern`")
	}
	if err := regexEnricher.TryAdd(devops.PRODUCTION, op.ProductionPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid value for `productionPattern`")
	}
	taskData := &tasks.BitbucketTaskData{
		Options:       op,
		ApiClient:     apiClient,
		RegexEnricher: regexEnricher,
	}
	if !timeAfter.IsZero() {
		taskData.TimeAfter = &timeAfter
		logger.Debug("collect data updated timeAfter %s", timeAfter)
	}

	return taskData, nil
}

func (p Bitbucket) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/bitbucket"
}

func (p Bitbucket) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Bitbucket) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (pp plugin.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (p Bitbucket) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/scopes/*scopeId": {
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
	}
}

func (p Bitbucket) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.BitbucketTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}

func EnrichOptions(taskCtx plugin.TaskContext,
	op *tasks.BitbucketOptions,
	apiClient *helper.ApiClient) errors.Error {
	var repo models.BitbucketRepo
	// validate the op and set name=owner/repo if this is from advanced mode or bpV100
	err := tasks.ValidateTaskOptions(op)
	if err != nil {
		return err
	}
	logger := taskCtx.GetLogger()
	err = taskCtx.GetDal().First(&repo, dal.Where(
		"connection_id = ? AND bitbucket_id = ?",
		op.ConnectionId, op.FullName))
	if err == nil {
		if op.TransformationRuleId == 0 {
			op.TransformationRuleId = repo.TransformationRuleId
		}
	} else {
		if taskCtx.GetDal().IsErrorNotFound(err) && op.FullName != "" {
			var repo *models.BitbucketApiRepo
			repo, err = tasks.GetApiRepo(op, apiClient)
			if err != nil {
				return err
			}
			logger.Debug(fmt.Sprintf("Current repo: %s", repo.FullName))
			scope := repo.ConvertApiScope().(*models.BitbucketRepo)
			scope.ConnectionId = op.ConnectionId
			err = taskCtx.GetDal().CreateIfNotExist(scope)
			if err != nil {
				return err
			}
		} else {
			return errors.Default.Wrap(err, fmt.Sprintf("fail to find repo %s", op.FullName))
		}
	}
	// Set GithubTransformationRule if it's nil, this has lower priority
	if op.BitbucketTransformationRule == nil && op.TransformationRuleId != 0 {
		var transformationRule models.BitbucketTransformationRule
		db := taskCtx.GetDal()
		err = db.First(&transformationRule, dal.Where("id = ?", repo.TransformationRuleId))
		if err != nil && !db.IsErrorNotFound(err) {
			return errors.BadInput.Wrap(err, "fail to get transformationRule")
		}
		op.BitbucketTransformationRule = &transformationRule
	}
	if op.BitbucketTransformationRule == nil && op.TransformationRuleId == 0 {
		op.BitbucketTransformationRule = new(models.BitbucketTransformationRule)
	}
	return err
}
