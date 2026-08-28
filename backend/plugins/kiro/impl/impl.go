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
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/kiro/api"
	"github.com/apache/incubator-devlake/plugins/kiro/models"
	"github.com/apache/incubator-devlake/plugins/kiro/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/kiro/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginSource
	plugin.PluginMigration
	plugin.DataSourcePluginBlueprintV200
} = (*Kiro)(nil)

// Kiro collects Kiro enterprise usage exports from S3.
//
// This is a separate plugin rather than an evolution of the retired predecessor:
// the old format is frozen, and a plugin that keeps evolving should not depend
// on frozen code. No implementation code is shared, following the same split
// as bitbucket and bitbucket_server.
type Kiro struct{}

func (p Kiro) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)
	return nil
}

func (p Kiro) Name() string {
	return "kiro"
}

func (p Kiro) Description() string {
	return "collect Kiro usage reports and interaction logs from S3"
}

func (p Kiro) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/kiro"
}

// GetTablesInfo must list every model or plugins/table_info_test.go fails.
func (p Kiro) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.KiroConnection{},
		&models.KiroS3Slice{},
		&models.KiroS3FileMeta{},
		&models.KiroUserReport{},
		&models.KiroUserModelMessage{},
		&models.KiroChatLog{},
		&models.KiroCompletionLog{},
	}
}

func (p Kiro) Connection() dal.Tabler {
	return &models.KiroConnection{}
}

func (p Kiro) Scope() plugin.ToolLayerScope {
	return &models.KiroS3Slice{}
}

// ScopeConfig returns nil: the export format is defined by AWS and uniform
// across an organization, so there is nothing per-scope to configure.
func (p Kiro) ScopeConfig() dal.Tabler {
	return nil
}

func (p Kiro) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

// SubTaskMetas lists discovery first, then one extractor per stream. The
// extractors declare their dependency on discovery, so the split is safe and
// gives each stream its own progress reporting - useful when a scope holds tens
// of thousands of log objects and one needs to know which stream is slow.
func (p Kiro) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectKiroS3FilesMeta,
		tasks.ExtractKiroUserReportMeta,
		tasks.ExtractKiroChatLogMeta,
		tasks.ExtractKiroCompletionLogMeta,
	}
}

func (p Kiro) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.KiroOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is required")
	}
	if op.AccountId == "" {
		return nil, errors.BadInput.New("accountId is required")
	}
	if op.Year <= 0 {
		return nil, errors.BadInput.New("year is required")
	}

	connectionHelper := helper.NewConnectionHelper(taskCtx, nil, p.Name())
	connection := &models.KiroConnection{}
	if err := connectionHelper.FirstById(connection, op.ConnectionId); err != nil {
		return nil, err
	}

	s3Clients, err := tasks.NewKiroS3Clients(connection)
	if err != nil {
		return nil, err
	}

	// Identity Store is optional and only supplies display names, so a failure
	// here degrades presentation rather than collection.
	identityClient, identityErr := tasks.NewKiroIdentityClient(connection)
	if identityErr != nil {
		taskCtx.GetLogger().Warn(identityErr, "identity store unavailable, proceeding without display names")
		identityClient = nil
	}

	timePath := fmt.Sprintf("%04d", op.Year)
	if op.Month != nil {
		timePath = fmt.Sprintf("%04d/%02d", op.Year, *op.Month)
	}

	return &tasks.KiroTaskData{
		Options:        &op,
		Connection:     connection,
		S3Clients:      s3Clients,
		IdentityClient: identityClient,
		Prefixes:       tasks.BuildPrefixes(connection, op.AccountId, timePath),
	}, nil
}

func (p Kiro) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes)
}

func (p Kiro) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
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
		"connections/:connectionId/test": {
			"POST": api.TestExistingConnection,
		},
		// Scope discovery: lists the accounts, years and months that actually
		// have exported data, so a scope is selected instead of hand-entered.
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": api.SearchRemoteScopes,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScopes,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    api.GetScope,
			"PATCH":  api.PatchScope,
			"DELETE": api.DeleteScope,
		},
		"connections/:connectionId/scopes/:scopeId/latest-sync-state": {
			"GET": api.GetScopeLatestSyncState,
		},
	}
}
