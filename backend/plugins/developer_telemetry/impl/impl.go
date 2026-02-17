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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/developer_telemetry/api"
	"github.com/apache/incubator-devlake/plugins/developer_telemetry/models"
	"github.com/apache/incubator-devlake/plugins/developer_telemetry/models/migrationscripts"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginModel
	plugin.PluginMigration
	plugin.PluginApi
} = (*DeveloperTelemetry)(nil)

type DeveloperTelemetry struct{}

func (p DeveloperTelemetry) Connection() dal.Tabler {
	return &models.DeveloperTelemetryConnection{}
}

func (p DeveloperTelemetry) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)
	return nil
}

func (p DeveloperTelemetry) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.DeveloperTelemetryConnection{},
		&models.DeveloperMetrics{},
	}
}

func (p DeveloperTelemetry) Description() string {
	return "Collect developer telemetry data from local development environments"
}

func (p DeveloperTelemetry) Name() string {
	return "developer_telemetry"
}

func (p DeveloperTelemetry) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/developer_telemetry"
}

func (p DeveloperTelemetry) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p DeveloperTelemetry) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
		"connections/:connectionId/test": {
			"POST": api.TestConnection,
		},
		"connections/:connectionId/report": {
			"POST": api.PostReport,
		},
	}
}
