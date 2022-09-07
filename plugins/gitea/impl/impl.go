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

	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitea/api"
	"github.com/apache/incubator-devlake/plugins/gitea/models"
	"github.com/apache/incubator-devlake/plugins/gitea/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/gitea/tasks"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Gitea)(nil)
var _ core.PluginInit = (*Gitea)(nil)
var _ core.PluginTask = (*Gitea)(nil)
var _ core.PluginApi = (*Gitea)(nil)
var _ core.Migratable = (*Gitea)(nil)
var _ core.CloseablePluginTask = (*Gitea)(nil)

type Gitea string

func (plugin Gitea) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	api.Init(config, logger, db)
	return nil
}

func (plugin Gitea) GetTablesInfo() []core.Tabler {
	return []core.Tabler{
		&models.GiteaConnection{},
		&models.GiteaAccount{},
		&models.GiteaCommit{},
		&models.GiteaCommitStat{},
		&models.GiteaIssue{},
		&models.GiteaIssueComment{},
		&models.GiteaIssueLabel{},
		&models.GiteaRepo{},
		&models.GiteaRepoCommit{},
		&models.GiteaResponse{},
		&models.GiteaReviewer{},
	}
}

func (plugin Gitea) Description() string {
	return "To collect and enrich data from Gitea"
}

func (plugin Gitea) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectApiRepoMeta,
		tasks.ExtractApiRepoMeta,
		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,
		tasks.CollectCommitsMeta,
		tasks.ExtractCommitsMeta,
		tasks.CollectApiIssueCommentsMeta,
		tasks.ExtractApiIssueCommentsMeta,
		tasks.CollectApiCommitStatsMeta,
		tasks.ExtractApiCommitStatsMeta,
		tasks.ConvertRepoMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertCommitsMeta,
		tasks.ConvertIssueLabelsMeta,
		tasks.ConvertAccountsMeta,
		tasks.ConvertIssueCommentsMeta,
	}
}

func (plugin Gitea) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.GiteaOptions
	var err error
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}

	if op.Owner == "" {
		return nil, fmt.Errorf("owner is required for Gitea execution")
	}

	if op.Repo == "" {
		return nil, fmt.Errorf("repo is required for Gitea execution")
	}

	if op.ConnectionId == 0 {
		return nil, fmt.Errorf("connectionId is invalid")
	}

	connection := &models.GiteaConnection{}
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
	apiClient, err := tasks.NewGiteaApiClient(taskCtx, connection)

	if err != nil {
		return nil, err
	}

	return &tasks.GiteaTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin Gitea) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gitea"
}

func (plugin Gitea) MigrationScripts() []migration.Script {
	return migrationscripts.All()
}

func (plugin Gitea) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
	}
}

func (plugin Gitea) MakePipelinePlan(connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, error) {
	return api.MakePipelinePlan(plugin.SubTaskMetas(), connectionId, scope)
}

func (plugin Gitea) Close(taskCtx core.TaskContext) error {
	data, ok := taskCtx.GetData().(*tasks.GiteaTaskData)
	if !ok {
		return fmt.Errorf("GetData failed when try to close %+v", taskCtx)
	}
	data.ApiClient.Release()
	return nil
}
