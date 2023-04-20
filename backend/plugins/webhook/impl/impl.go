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
	"github.com/apache/incubator-devlake/plugins/webhook/api"
	"github.com/apache/incubator-devlake/plugins/webhook/models/migrationscripts"
)

// make sure interface is implemented
var _ plugin.PluginMeta = (*Webhook)(nil)
var _ plugin.PluginInit = (*Webhook)(nil)
var _ plugin.PluginApi = (*Webhook)(nil)
var _ plugin.PluginModel = (*Webhook)(nil)
var _ plugin.PluginMigration = (*Webhook)(nil)

type Webhook struct{}

func (p Webhook) Description() string {
	return "collect some Webhook data"
}

func (p Webhook) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (p Webhook) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{}
}

func (p Webhook) MakeDataSourcePipelinePlanV200(connectionId uint64, _ []*plugin.BlueprintScopeV200, _ plugin.BlueprintSyncPolicy) (pp plugin.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(connectionId)
}

// PkgPath information lost when compiled as plugin(.so)
func (p Webhook) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/webhook"
}

func (p Webhook) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Webhook) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
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
		":connectionId/deployments": {
			"POST": api.PostDeploymentCicdTask,
		},
		":connectionId/issues": {
			"POST": api.PostIssue,
		},
		":connectionId/issue/:issueKey/close": {
			"POST": api.CloseIssue,
		},
	}
}
