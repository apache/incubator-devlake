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
	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/feishu/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/feishu/tasks"
	"github.com/apache/incubator-devlake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Feishu)(nil)
var _ core.PluginInit = (*Feishu)(nil)
var _ core.PluginTask = (*Feishu)(nil)
var _ core.PluginApi = (*Feishu)(nil)
var _ core.Migratable = (*Feishu)(nil)

type Feishu struct{}

func (plugin Feishu) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	return nil
}

func (plugin Feishu) Description() string {
	return "To collect and enrich data from Feishu"
}

func (plugin Feishu) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectMeetingTopUserItemMeta,
		tasks.ExtractMeetingTopUserItemMeta,
	}
}

func (plugin Feishu) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.FeishuOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	apiClient, err := tasks.NewFeishuApiClient(taskCtx)
	if err != nil {
		return nil, err
	}
	return &tasks.FeishuTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin Feishu) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/feishu"
}

func (plugin Feishu) MigrationScripts() []migration.Script {
	return []migration.Script{
		new(migrationscripts.InitSchemas), new(migrationscripts.UpdateSchemas20220524),
	}
}

func (plugin Feishu) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{}
}

var PluginEntry Feishu

// standalone mode for debugging
func main() {
	feishuCmd := &cobra.Command{Use: "feishu"}
	numOfDaysToCollect := feishuCmd.Flags().IntP("numOfDaysToCollect", "n", 8, "feishu collect days")
	_ = feishuCmd.MarkFlagRequired("numOfDaysToCollect")
	feishuCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"numOfDaysToCollect": *numOfDaysToCollect,
		})
	}
	runner.RunCmd(feishuCmd)
}
