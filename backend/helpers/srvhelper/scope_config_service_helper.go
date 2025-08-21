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
	"fmt"
	"strings"

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
	s := new(S)
	// find out the primary key of the scope model
	sPk := errors.Must1(dal.GetPrimarykeyColumnNames(scopeConfigSrv.db, (interface{}(s)).(dal.Tabler)))
	if len(sPk) != 2 {
		return nil, errors.Internal.New("Scope model should have 2 primary key fields")
	}
	theOtherPk := sPk[0]
	if strings.HasSuffix(theOtherPk, ".connection_id") {
		theOtherPk = sPk[1]
	}
	var bpss []struct {
		S           *S `gorm:"embedded"`
		BlueprintId uint64
		ProjectName string
		ScopeId     string
	}
	scopeTable := (*s).TableName()
	// Postgres fails as scope_id is a varchar and theOtherPk can be an integer in some cases
	join := fmt.Sprintf("LEFT JOIN %s ON (%s.connection_id = bps.connection_id AND %s = bps.scope_id)", scopeTable, scopeTable, theOtherPk)
	if scopeConfigSrv.db.Dialect() == "postgres" {
		join = fmt.Sprintf("LEFT JOIN %s ON (%s.connection_id = bps.connection_id AND CAST(%s AS varchar) = bps.scope_id)", scopeTable, scopeTable, theOtherPk)
	}
	errors.Must(scopeConfigSrv.db.All(
		&bpss,
		dal.Select(fmt.Sprintf("bp.id AS blueprint_id, bp.project_name, bps.scope_id, %s.*", scopeTable)),
		dal.From("_devlake_blueprint_scopes bps"),
		dal.Join("LEFT JOIN _devlake_blueprints bp ON (bp.id = bps.blueprint_id)"),
		dal.Join(join),
		dal.Where("bps.plugin_name = ? AND bps.connection_id = ? AND scope_config_id = ?", pluginName, (*scopeConfig).ScopeConfigConnectionId(), (*scopeConfig).ScopeConfigId()),
	))
	projectScopeMap := make(map[string]*models.ProjectScope)
	for _, bps := range bpss {
		if _, ok := projectScopeMap[bps.ProjectName]; !ok {
			projectScopeMap[bps.ProjectName] = &models.ProjectScope{
				Name:        bps.ProjectName,
				BlueprintId: bps.BlueprintId,
			}
		}
		projectScopeMap[bps.ProjectName].Scopes = append(
			projectScopeMap[bps.ProjectName].Scopes,
			struct {
				ScopeID   string `json:"scopeId"`
				ScopeName string `json:"scopeName"`
			}{
				ScopeID:   bps.ScopeId,
				ScopeName: (*bps.S).ScopeName(),
			},
		)
	}
	ps := &models.ProjectScopeOutput{}
	for _, projectScope := range projectScopeMap {
		ps.Projects = append(ps.Projects, *projectScope)
	}
	ps.Count = len(ps.Projects)
	return ps, nil
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
