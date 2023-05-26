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
	"github.com/apache/incubator-devlake/plugins/tiktokAds/api"
	"github.com/apache/incubator-devlake/plugins/tiktokAds/models"
	"github.com/apache/incubator-devlake/plugins/tiktokAds/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/tiktokAds/tasks"
)

var _ plugin.PluginMeta = (*TiktokAds)(nil)
var _ plugin.PluginInit = (*TiktokAds)(nil)
var _ plugin.PluginTask = (*TiktokAds)(nil)
var _ plugin.PluginApi = (*TiktokAds)(nil)
var _ plugin.PluginModel = (*TiktokAds)(nil)
var _ plugin.PluginMigration = (*TiktokAds)(nil)
var _ plugin.CloseablePluginTask = (*TiktokAds)(nil)

type TiktokAds struct{}

func (p TiktokAds) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p TiktokAds) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.TiktokAdsConnection{},
		&models.TiktokAdsMeetingTopUserItem{},
	}
}

func (p TiktokAds) Description() string {
	return "To collect and enrich data from TiktokAds"
}

func (p TiktokAds) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectMeetingTopUserItemMeta,
		tasks.ExtractMeetingTopUserItemMeta,
	}
}

func (p TiktokAds) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.TiktokAdsOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}

	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.TiktokAdsConnection{}
	err := connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.NewTiktokAdsApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}
	return &tasks.TiktokAdsTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (p TiktokAds) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/tiktokAds"
}

func (p TiktokAds) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p TiktokAds) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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

func (p TiktokAds) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.TiktokAdsTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
