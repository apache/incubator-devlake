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

package srvhelper

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
)

// NoScopeConfig is a placeholder for plugins that don't have any scope configuration yet
type NoScopeConfig struct{}

func (NoScopeConfig) TableName() string               { return "" }
func (NoScopeConfig) ScopeConfigId() uint64           { return 0 }
func (NoScopeConfig) ScopeConfigConnectionId() uint64 { return 0 }

// ScopeConfigSrvHelper
type ScopeConfigSrvHelper[C plugin.ToolLayerConnection, S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig] struct {
	*ModelSrvHelper[SC]
}

func NewScopeConfigSrvHelper[
	C plugin.ToolLayerConnection,
	S plugin.ToolLayerScope,
	SC plugin.ToolLayerScopeConfig,
](basicRes context.BasicRes, searchColumns []string) *ScopeConfigSrvHelper[C, S, SC] {
	return &ScopeConfigSrvHelper[C, S, SC]{
		ModelSrvHelper: NewModelSrvHelper[SC](basicRes, searchColumns),
	}
}

func (scopeConfigSrv *ScopeConfigSrvHelper[C, S, SC]) GetAllByConnectionId(connectionId uint64) ([]*SC, errors.Error) {
	var scopeConfigs []*SC
	err := scopeConfigSrv.db.All(&scopeConfigs,
		dal.Where("connection_id = ?", connectionId),
		dal.Orderby("id DESC"),
	)
	return scopeConfigs, err
}

func (scopeConfigSrv *ScopeConfigSrvHelper[C, S, SC]) GetProjectsByScopeConfig(pluginName string, scopeConfig *SC) (*models.ProjectScopeOutput, errors.Error) {
	ps := &models.ProjectScopeOutput{}
	projectMap := make(map[string]*models.ProjectScope)
	// 1. get all scopes that are using the scopeConfigId
	var scope []*S
	err := scopeConfigSrv.db.All(&scope,
		dal.Where("scope_config_id = ?", (*scopeConfig).ScopeConfigId()),
	)
	if err != nil {
		return nil, err
	}
	for _, s := range scope {
		// 2. get blueprint id by connection id and scope id
		bpScope := []*models.BlueprintScope{}
		err = scopeConfigSrv.db.All(&bpScope,
			dal.Where("plugin_name = ? and connection_id = ? and scope_id = ?", pluginName, (*s).ScopeConnectionId(), (*s).ScopeId()),
		)
		if err != nil {
			return nil, err
		}

		for _, bs := range bpScope {
			// 3. get project details by blueprint id
			bp := models.Blueprint{}
			err = scopeConfigSrv.db.All(&bp,
				dal.Where("id = ?", bs.BlueprintId),
			)
			if err != nil {
				return nil, err
			}
			if project, exists := projectMap[bp.ProjectName]; exists {
				project.Scopes = append(project.Scopes, struct {
					ScopeID   string `json:"scopeId"`
					ScopeName string `json:"scopeName"`
				}{
					ScopeID:   bs.ScopeId,
					ScopeName: (*s).ScopeName(),
				})
			} else {
				projectMap[bp.ProjectName] = &models.ProjectScope{
					Name:        bp.ProjectName,
					BlueprintId: bp.ID,
					Scopes: []struct {
						ScopeID   string `json:"scopeId"`
						ScopeName string `json:"scopeName"`
					}{
						{
							ScopeID:   bs.ScopeId,
							ScopeName: (*s).ScopeName(),
						},
					},
				}
			}
		}
	}
	// 4. combine all projects
	for _, project := range projectMap {
		ps.Projects = append(ps.Projects, *project)
	}
	ps.Count = len(ps.Projects)

	return ps, err
}

func (scopeConfigSrv *ScopeConfigSrvHelper[C, S, SC]) DeleteScopeConfig(scopeConfig *SC) (refs []*S, err errors.Error) {
	err = scopeConfigSrv.ModelSrvHelper.NoRunningPipeline(func(tx dal.Transaction) errors.Error {
		// make sure no scope is using the scopeConfig
		sc := *scopeConfig
		errors.Must(tx.All(
			&refs,
			dal.Where("connection_id = ? AND scope_config_id = ?", sc.ScopeConfigConnectionId(), sc.ScopeConfigId()),
		))
		if len(refs) > 0 {
			return errors.Conflict.New("Please delete all data scope(s) before you delete this ScopeConfig.")
		}
		errors.Must(tx.Delete(scopeConfig))
		return nil
	})
	return
}
