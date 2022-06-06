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

package main

import (
	"fmt"

	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/jenkins/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"
	"github.com/apache/incubator-devlake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Jenkins)(nil)
var _ core.PluginInit = (*Jenkins)(nil)
var _ core.PluginTask = (*Jenkins)(nil)
var _ core.PluginApi = (*Jenkins)(nil)
var _ core.Migratable = (*Jenkins)(nil)

type Jenkins struct{}

func (plugin Jenkins) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	return nil
}

func (plugin Jenkins) Description() string {
	return "To collect and enrich data from Jenkins"
}

func (plugin Jenkins) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectApiJobsMeta,
		tasks.ExtractApiJobsMeta,
		tasks.CollectApiBuildsMeta,
		tasks.ExtractApiBuildsMeta,
		tasks.ConvertJobsMeta,
		tasks.ConvertBuildsMeta,
	}
}
func (plugin Jenkins) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.JenkinsOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode options: %v", err)
	}
	apiClient, err := tasks.CreateApiClient(taskCtx)
	if err != nil {
		return nil, err
	}
	return &tasks.JenkinsTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin Jenkins) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/jenkins"
}

func (plugin Jenkins) MigrationScripts() []migration.Script {
	return []migration.Script{new(migrationscripts.InitSchemas)}
}

func (plugin Jenkins) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"GET": api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":   api.GetConnection,
			"PATCH": api.PatchConnection,
		},
	}
}

var PluginEntry Jenkins //nolint

func main() {
	jenkinsCmd := &cobra.Command{Use: "jenkins"}
	jenkinsCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, []string{}, map[string]interface{}{})
	}
	runner.RunCmd(jenkinsCmd)
}
