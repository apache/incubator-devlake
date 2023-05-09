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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"
)

var _ plugin.PluginMeta = (*Jenkins)(nil)
var _ plugin.PluginInit = (*Jenkins)(nil)
var _ plugin.PluginTask = (*Jenkins)(nil)
var _ plugin.PluginApi = (*Jenkins)(nil)
var _ plugin.PluginModel = (*Jenkins)(nil)
var _ plugin.PluginMigration = (*Jenkins)(nil)
var _ plugin.CloseablePluginTask = (*Jenkins)(nil)
var _ plugin.PluginSource = (*Jenkins)(nil)

type Jenkins struct{}

func (p Jenkins) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p Jenkins) Connection() interface{} {
	return &models.JenkinsConnection{}
}

func (p Jenkins) Scope() interface{} {
	return &models.JenkinsJob{}
}

func (p Jenkins) TransformationRule() interface{} {
	return &models.JenkinsTransformationRule{}
}

func (p Jenkins) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.JenkinsBuild{},
		&models.JenkinsBuildCommit{},
		&models.JenkinsConnection{},
		&models.JenkinsJob{},
		&models.JenkinsJobDag{},
		&models.JenkinsPipeline{},
		&models.JenkinsStage{},
		&models.JenkinsTask{},
	}
}

func (p Jenkins) Description() string {
	return "To collect and enrich data from Jenkins"
}

func (p Jenkins) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.ConvertJobsMeta,
		tasks.CollectApiBuildsMeta,
		tasks.ExtractApiBuildsMeta,
		tasks.CollectApiStagesMeta,
		tasks.ExtractApiStagesMeta,
		tasks.EnrichApiBuildWithStagesMeta,
		tasks.ConvertBuildsToCICDMeta,
		tasks.ConvertStagesMeta,
		tasks.ConvertBuildReposMeta,
	}
}
func (p Jenkins) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	// Firstly, let's decode options to JenkinsOptions
	op, err := tasks.DecodeTaskOptions(options)
	if err != nil {
		return nil, err
	}

	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	connection := &models.JenkinsConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	err = EnrichOptions(taskCtx, op, apiClient)
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
	taskData := &tasks.JenkinsTaskData{
		Options:       op,
		ApiClient:     apiClient,
		Connection:    connection,
		RegexEnricher: regexEnricher,
	}
	if !timeAfter.IsZero() {
		taskData.TimeAfter = &timeAfter
		logger.Debug("collect data created from %s", timeAfter)
	}
	return taskData, nil
}

func (p Jenkins) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/jenkins"
}

func (p Jenkins) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Jenkins) MakePipelinePlan(connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	return api.MakePipelinePlanV100(p.SubTaskMetas(), connectionId, scope)
}

func (p Jenkins) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (pp plugin.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (p Jenkins) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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

func (p Jenkins) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.JenkinsTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}

func EnrichOptions(taskCtx plugin.TaskContext,
	op *tasks.JenkinsOptions,
	apiClient *helper.ApiAsyncClient) errors.Error {
	jenkinsJob := &models.JenkinsJob{}
	// If this is from BpV200, we should set JobFullName to scopeId
	if op.JobFullName == "" {
		op.JobFullName = op.ScopeId
	}
	// validate the op and set name=owner/repo if this is from advanced mode or bpV100
	op, err := tasks.ValidateTaskOptions(op)
	if err != nil {
		return err
	}
	log := taskCtx.GetLogger()

	// for advanced mode or others which we only have name, for bp v200, we have TransformationRuleId
	err = taskCtx.GetDal().First(jenkinsJob,
		dal.Where(`connection_id = ? and full_name = ?`,
			op.ConnectionId, op.JobFullName))
	if err == nil {
		if op.TransformationRuleId == 0 {
			op.TransformationRuleId = jenkinsJob.TransformationRuleId
		}
	}

	err = api.GetJob(apiClient, op.JobPath, op.JobName, op.JobFullName, 100, func(job *models.Job, isPath bool) errors.Error {
		log.Debug(fmt.Sprintf("Current job: %s", job.FullName))
		op.Name = job.Name
		op.JobPath = job.Path
		jenkinsJob := ConvertJobToJenkinsJob(job, op)
		err = taskCtx.GetDal().CreateIfNotExist(jenkinsJob)
		return err
	})
	if err != nil {
		return err
	}

	if !strings.HasSuffix(op.JobPath, "/") {
		op.JobPath = fmt.Sprintf("%s/", op.JobPath)
	}
	// We only set op.JenkinsTransformationRule when it's nil and we have op.TransformationRuleId != 0
	if op.JenkinsTransformationRule.DeploymentPattern == "" && op.JenkinsTransformationRule.ProductionPattern == "" && op.TransformationRuleId != 0 {
		var transformationRule models.JenkinsTransformationRule
		err = taskCtx.GetDal().First(&transformationRule, dal.Where("id = ?", op.TransformationRuleId))
		if err != nil {
			return errors.BadInput.Wrap(err, "fail to get transformationRule")
		}
		op.JenkinsTransformationRule = &transformationRule
	}

	if op.JenkinsTransformationRule.DeploymentPattern == "" && op.JenkinsTransformationRule.ProductionPattern == "" && op.TransformationRuleId == 0 {
		op.JenkinsTransformationRule = new(models.JenkinsTransformationRule)
	}

	return nil
}

func ConvertJobToJenkinsJob(job *models.Job, op *tasks.JenkinsOptions) *models.JenkinsJob {
	return &models.JenkinsJob{
		ConnectionId:         op.ConnectionId,
		FullName:             job.FullName,
		TransformationRuleId: op.TransformationRuleId,
		Name:                 job.Name,
		Path:                 job.Path,
		Class:                job.Class,
		Color:                job.Color,
		Base:                 job.Base,
		Url:                  job.URL,
		Description:          job.Description,
		PrimaryView:          job.URL + job.Path + job.Class,
	}
}
