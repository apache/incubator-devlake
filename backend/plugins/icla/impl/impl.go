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
	"github.com/apache/incubator-devlake/plugins/icla/models"
	"github.com/apache/incubator-devlake/plugins/icla/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/icla/tasks"
)

// make sure interface is implemented
var _ plugin.PluginMeta = (*Icla)(nil)
var _ plugin.PluginInit = (*Icla)(nil)
var _ plugin.PluginTask = (*Icla)(nil)
var _ plugin.PluginApi = (*Icla)(nil)
var _ plugin.PluginModel = (*Icla)(nil)
var _ plugin.PluginMigration = (*Icla)(nil)
var _ plugin.CloseablePluginTask = (*Icla)(nil)

type Icla struct{}

func (p Icla) Description() string {
	return "collect some Icla data"
}

func (p Icla) Init(basicRes context.BasicRes) errors.Error {
	return nil
}

func (p Icla) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.IclaCommitter{},
	}
}

func (p Icla) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectCommitterMeta,
		tasks.ExtractCommitterMeta,
	}
}

func (p Icla) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.IclaOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}

	apiClient, err := errors.Convert01(tasks.NewIclaApiClient(taskCtx))
	if err != nil {
		return nil, err
	}

	return &tasks.IclaTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (p Icla) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/icla"
}

func (p Icla) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Icla) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return nil
}

func (p Icla) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.IclaTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
