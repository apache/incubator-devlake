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
	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Gitlab)(nil)
var _ core.PluginInit = (*Gitlab)(nil)
var _ core.PluginTask = (*Gitlab)(nil)
var _ core.PluginApi = (*Gitlab)(nil)
var _ core.Migratable = (*Gitlab)(nil)

type Gitlab string

func (plugin Gitlab) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	return nil
}

func (plugin Gitlab) Description() string {
	return "To collect and enrich data from Gitlab"
}

func (plugin Gitlab) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectProjectMeta,
		tasks.ExtractProjectMeta,
		tasks.CollectCommitsMeta,
		tasks.ExtractCommitsMeta,
		tasks.CollectTagMeta,
		tasks.ExtractTagMeta,
		tasks.CollectApiMergeRequestsMeta,
		tasks.ExtractApiMergeRequestsMeta,
		tasks.CollectApiMergeRequestsNotesMeta,
		tasks.ExtractApiMergeRequestsNotesMeta,
		tasks.CollectApiMergeRequestsCommitsMeta,
		tasks.ExtractApiMergeRequestsCommitsMeta,
		tasks.CollectApiPipelinesMeta,
		tasks.ExtractApiPipelinesMeta,
		tasks.CollectApiChildrenOnPipelinesMeta,
		tasks.ExtractApiChildrenOnPipelinesMeta,
		tasks.EnrichMergeRequestsMeta,
		tasks.ConvertProjectMeta,
		tasks.ConvertApiMergeRequestsMeta,
		tasks.ConvertApiCommitsMeta,
		tasks.ConvertApiNotesMeta,
		tasks.ConvertMergeRequestCommentMeta,
		tasks.ConvertApiMergeRequestsCommitsMeta,
	}
}

func (plugin Gitlab) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.GitlabOptions
	var err error
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.NewGitlabApiClient(taskCtx)
	if err != nil {
		return nil, err
	}

	return &tasks.GitlabTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin Gitlab) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gitlab"
}

func (plugin Gitlab) MigrationScripts() []migration.Script {
	return []migration.Script{new(migrationscripts.InitSchemas), new(migrationscripts.UpdateSchemas20220510)}
}

func (plugin Gitlab) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"GET": api.ListConnections,
		},
		"connections/:connectionId": {
			"GET": api.GetConnection,
			"PATCH": api.PatchConnection,
		},
	}
}
