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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
	"time"
)

var _ plugin.PluginMeta = (*Bitbucket)(nil)
var _ plugin.PluginInit = (*Bitbucket)(nil)
var _ plugin.PluginTask = (*Bitbucket)(nil)
var _ plugin.PluginApi = (*Bitbucket)(nil)
var _ plugin.PluginModel = (*Bitbucket)(nil)
var _ plugin.PluginMigration = (*Bitbucket)(nil)
var _ plugin.PluginBlueprintV100 = (*Bitbucket)(nil)
var _ plugin.CloseablePluginTask = (*Bitbucket)(nil)

type Bitbucket string

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
		tasks.CollectApiRepoMeta,
		tasks.ExtractApiRepoMeta,

		tasks.CollectApiPullRequestsMeta,
		tasks.ExtractApiPullRequestsMeta,

		tasks.CollectApiPrCommentsMeta,
		tasks.ExtractApiPrCommentsMeta,

		tasks.CollectApiPrCommitsMeta,
		tasks.ExtractApiPrCommitsMeta,

		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,

		tasks.CollectApiIssueCommentsMeta,
		tasks.ExtractApiIssueCommentsMeta,

		tasks.CollectApiPipelinesMeta,
		tasks.ExtractApiPipelinesMeta,

		tasks.CollectApiDeploymentsMeta,
		tasks.ExtractApiDeploymentsMeta,

		tasks.ConvertRepoMeta,
		tasks.ConvertAccountsMeta,
		tasks.ConvertPullRequestsMeta,
		tasks.ConvertPrCommentsMeta,
		tasks.ConvertPrCommitsMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertIssueCommentsMeta,
		tasks.ConvertPipelineMeta,
		tasks.ConvertDeploymentMeta,
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

	var createdDateAfter time.Time
	if op.CreatedDateAfter != "" {
		createdDateAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.CreatedDateAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `createdDateAfter`")
		}
	}
	taskData := &tasks.BitbucketTaskData{
		Options:   op,
		ApiClient: apiClient,
	}
	if !createdDateAfter.IsZero() {
		taskData.CreatedDateAfter = &createdDateAfter
		logger.Debug("collect data updated createdDateAfter %s", createdDateAfter)
	}

	return taskData, nil
}

func (p Bitbucket) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/bitbucket"
}

func (p Bitbucket) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Bitbucket) MakePipelinePlan(connectionId uint64, scope []*plugin.BlueprintScopeV100) (plugin.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(p.SubTaskMetas(), connectionId, scope)
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
