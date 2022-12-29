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

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"
)

var _ core.PluginMeta = (*Jenkins)(nil)
var _ core.PluginInit = (*Jenkins)(nil)
var _ core.PluginTask = (*Jenkins)(nil)
var _ core.PluginApi = (*Jenkins)(nil)
var _ core.PluginModel = (*Jenkins)(nil)
var _ core.PluginMigration = (*Jenkins)(nil)
var _ core.CloseablePluginTask = (*Jenkins)(nil)
var _ core.PluginSource = (*Jenkins)(nil)

type Jenkins struct{}

func (plugin Jenkins) Init(basicRes core.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (plugin Jenkins) Connection() interface{} {
	return &models.JenkinsConnection{}
}

func (plugin Jenkins) Scope() interface{} {
	return &models.JenkinsJob{}
}

func (plugin Jenkins) TransformationRule() interface{} {
	return &models.JenkinsTransformationRule{}
}

func (plugin Jenkins) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.JenkinsBuild{},
		&models.JenkinsBuildCommit{},
		&models.JenkinsConnection{},
		&models.JenkinsJob{},
		&models.JenkinsJobDag{},
		&models.JenkinsPipeline{},
		&models.JenkinsResponse{},
		&models.JenkinsStage{},
		&models.JenkinsTask{},
	}
}

func (plugin Jenkins) Description() string {
	return "To collect and enrich data from Jenkins"
}

func (plugin Jenkins) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
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
func (plugin Jenkins) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
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

	var createdDateAfter time.Time
	if op.CreatedDateAfter != "" {
		createdDateAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.CreatedDateAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `createdDateAfter`")
		}
	}
	taskData := &tasks.JenkinsTaskData{
		Options:    op,
		ApiClient:  apiClient,
		Connection: connection,
	}
	if !createdDateAfter.IsZero() {
		taskData.CreatedDateAfter = &createdDateAfter
		logger.Debug("collect data created from %s", createdDateAfter)
	}
	return taskData, nil
}

func (plugin Jenkins) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/jenkins"
}

func (plugin Jenkins) MigrationScripts() []core.MigrationScript {
	return migrationscripts.All()
}

func (plugin Jenkins) MakePipelinePlan(connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, errors.Error) {
	return api.MakePipelinePlanV100(plugin.SubTaskMetas(), connectionId, scope)
}

func (plugin Jenkins) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*core.BlueprintScopeV200, syncPolicy core.BlueprintSyncPolicy) (pp core.PipelinePlan, sc []core.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(plugin.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (plugin Jenkins) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
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
		"connections/:connectionId/scopes/*fullName": {
			"GET":   api.GetScope,
			"PATCH": api.UpdateScope,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScope,
		},
		"transformation_rules": {
			"POST": api.CreateTransformationRule,
			"GET":  api.GetTransformationRuleList,
		},
		"transformation_rules/:id": {
			"PATCH": api.UpdateTransformationRule,
			"GET":   api.GetTransformationRule,
		},
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
	}
}

func (plugin Jenkins) Close(taskCtx core.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.JenkinsTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}

func EnrichOptions(taskCtx core.TaskContext,
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
		op.Name = jenkinsJob.Name
		op.JobPath = jenkinsJob.Path
		if op.TransformationRuleId == 0 {
			op.TransformationRuleId = jenkinsJob.TransformationRuleId
		}
	} else {
		if taskCtx.GetDal().IsErrorNotFound(err) && op.JobFullName != "" {

			pathSplit := strings.Split(op.JobFullName, "/")
			lastOne := len(pathSplit)

			path := "job/" + strings.Join(pathSplit[0:lastOne-1], "job/")
			if path == "job/" {
				path = ""
			}
			name := pathSplit[lastOne-1]

			err = api.GetJob(apiClient, path, name, op.JobFullName, 100, func(job *models.Job, isPath bool) errors.Error {
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
		} else {
			return errors.Default.Wrap(err, fmt.Sprintf("fail to find repo %s", op.Name))
		}
	}

	if !strings.HasSuffix(op.JobPath, "/") {
		op.JobPath = fmt.Sprintf("%s/", op.JobPath)
	}

	// We only set op.JenkinsTransformationRule when it's nil and we have op.TransformationRuleId != 0
	if op.JenkinsTransformationRule == nil && op.TransformationRuleId != 0 {
		var transformationRule models.JenkinsTransformationRule
		err = taskCtx.GetDal().First(&transformationRule, dal.Where("id = ?", op.TransformationRuleId))
		if err != nil {
			return errors.BadInput.Wrap(err, "fail to get transformationRule")
		}
		op.JenkinsTransformationRule = &transformationRule
	}

	if op.JenkinsTransformationRule == nil && op.TransformationRuleId == 0 {
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
