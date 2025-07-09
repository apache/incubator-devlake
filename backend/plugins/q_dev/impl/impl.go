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
	"github.com/apache/incubator-devlake/plugins/q_dev/api"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"github.com/apache/incubator-devlake/plugins/q_dev/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/q_dev/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginSource
	plugin.PluginMigration
	plugin.CloseablePluginTask
} = (*QDev)(nil)

type QDev struct{}

func (p QDev) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)
	return nil
}

func (p QDev) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.QDevConnection{},
		&models.QDevUserData{},
		&models.QDevS3FileMeta{},
	}
}

func (p QDev) Description() string {
	return "To collect and enrich data from AWS Q Developer usage metrics"
}

func (p QDev) Name() string {
	return "q_dev"
}

func (p QDev) Connection() dal.Tabler {
	return &models.QDevConnection{}
}

func (p QDev) Scope() plugin.ToolLayerScope {
	return nil
}

func (p QDev) ScopeConfig() dal.Tabler {
	return nil
}

func (p QDev) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectQDevS3FilesMeta,
		tasks.ExtractQDevS3DataMeta,
	}
}

func (p QDev) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.QDevOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}

	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	connection := &models.QDevConnection{}
	err := connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	// Create S3 client
	s3Client, err := tasks.NewQDevS3Client(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	// Create Identity client (new)
	identityClient, identityErr := tasks.NewQDevIdentityClient(connection)
	if identityErr != nil {
		taskCtx.GetLogger().Warn(identityErr, "Failed to create identity client, proceeding without user name resolution")
		identityClient = nil
	}

	return &tasks.QDevTaskData{
		Options:        &op,
		S3Client:       s3Client,
		IdentityClient: identityClient,
	}, nil
}

func (p QDev) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/q_dev"
}

func (p QDev) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p QDev) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		"connections/:connectionId/test": {
			"POST": api.TestExistingConnection,
		},
	}
}

func (p QDev) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.QDevTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.S3Client.Close()
	return nil
}
